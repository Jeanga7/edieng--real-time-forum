package api

import (
	"errors"
	"fmt"
	"forum/controllers"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch All Users
	var users []models.Model

	listOfPossibleParams := []string{"id", "email", "name", "username", "login", "firstname", "lastname", "password", "age", "gender"}

	nbParams := 0
	for _, v := range listOfPossibleParams {
		if len(r.FormValue(v)) > 0 {
			nbParams++
		}
	}
	paramName := ""
	expectedParams := true

	if nbParams == 1 {
		if len(r.FormValue("id")) > 0 {
			paramName = "id"
		} else if len(r.FormValue("email")) > 0 {
			paramName = "email"
		} else if len(r.FormValue("name")) > 0 {
			paramName = "name"
		} else if len(r.FormValue("username")) > 0 {
			paramName = "username"
		} else if len(r.FormValue("login")) > 0 {
			paramName = "login"
		} else {
			expectedParams = false
		}

		if expectedParams {
			if paramName == "login" {
				users = models.User.Fetch([]interface{}{"username"}, []interface{}{r.FormValue("login")})
				if len(users) == 0 {
					users = models.User.Fetch([]interface{}{"email"}, []interface{}{r.FormValue("login")})
				}
			} else {
				users = models.User.Fetch([]interface{}{paramName}, []interface{}{r.FormValue(paramName)})
			}
		} else {
			users = models.User.Fetch()
		}

	}

	fields := []string{"id", "name", "username", "email", "picture", "creationDate"}
	usersData := TransformModels(users, fields)
	WriteJSONResponse(w, usersData)
}

func ValidateFields(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)
	global.W = w

	form := html.EscapeString(r.FormValue("form"))
	if form == "login" {
		UsernameOrEmail := html.EscapeString(r.FormValue("login"))
		if len(UsernameOrEmail) > 0 {
			if strings.Contains(UsernameOrEmail, "@") && !helpers.ValidateInput("email", UsernameOrEmail) {
				helpers.JSONResponse(w, errors.New("INVALID EMAIL"), http.StatusBadRequest)
				return
			} else if !helpers.ValidateInput("username", UsernameOrEmail) {
				helpers.JSONResponse(w, errors.New("INVALID USERNAME"), http.StatusBadRequest)
				return
			}
		}

		password := html.EscapeString(r.FormValue("password"))
		if len(password) > 0 {
			if !helpers.ValidateInput("password", password) {
				helpers.JSONResponse(w, errors.New("INVALID PASSWORD"), http.StatusBadRequest)
				return
			}
		}

		user := models.User.Fetch([]interface{}{"email"}, []interface{}{UsernameOrEmail})
		if len(user) == 0 {
			user = models.User.Fetch([]interface{}{"username"}, []interface{}{UsernameOrEmail})
		}

		if len(user) == 1 {
			if len(password) > 0 {
				hashedPassword := fmt.Sprintf("%s", user[0].Values["password"])

				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
				if err == nil {
					controllers.CreateSession(w, r, fmt.Sprintf("%s", user[0].Values["id"]))
					controllers.BroadcastNewUser(user[0].Values["name"].(string))
					helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
					return
				} else {
					helpers.JSONResponse(w, errors.New("INVALID PASSWORD"), http.StatusNotFound)
					return
				}
			} else {
				helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
				return
			}
		} else if len(user) == 0 {
			helpers.JSONResponse(w, errors.New("USER NOT FOUND"), http.StatusNotFound)
			return
		}
	} else if form == "register" {
		name := html.EscapeString(r.FormValue("name"))
		if !helpers.ValidateInput("name", name) {
			helpers.JSONResponse(w, errors.New("INVALID NAME"), http.StatusBadRequest)
			return
		}

		email := html.EscapeString(r.FormValue("email"))
		if !helpers.ValidateInput("email", email) {
			helpers.JSONResponse(w, errors.New("INVALID EMAIL"), http.StatusBadRequest)
			return
		}

		username := html.EscapeString(r.FormValue("username"))
		if !helpers.ValidateInput("username", username) {
			helpers.JSONResponse(w, errors.New("INVALID USERNAME"), http.StatusBadRequest)
			return
		}

		password := html.EscapeString(r.FormValue("password"))
		if !helpers.ValidateInput("password", password) || len(password) < 8 || len(password) > 20 {
			helpers.JSONResponse(w, errors.New("INVALID PASSWORD"), http.StatusBadRequest)
			return
		}

		// Vérifiez si l'email ou le nom d'utilisateur est déjà utilisé
		resEmail := models.User.Fetch([]interface{}{"email"}, []interface{}{email})
		resUsername := models.User.Fetch([]interface{}{"username"}, []interface{}{username})
		if len(resEmail) > 0 || len(resUsername) > 0 {
			// Renvoyez une erreur si l'email ou le nom d'utilisateur est déjà utilisé
			helpers.JSONResponse(w, errors.New("USERNAME OR EMAIL ALREADY USED"), http.StatusConflict)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		helpers.ErrorSimple(err)

		picture, handler, _ := r.FormFile("picture")

		// Si non, créez un nouvel utilisateur
		newUser := models.User
		newUser.Values = map[string]interface{}{
			"name":     strings.ReplaceAll(strings.Title(name), " ", "_"),
			"username": username,
			"email":    email,
			"password": helpers.HashPassword(password),
		}

		if picture == nil {
			newUser.Values["picture"] = "default-user.png"
		} else {
			fileBytes, err := io.ReadAll(picture)
			if err != nil {
				helpers.ErrorSimple(err)
				helpers.JSONResponse(w, errors.New("FAILED TO READ IMAGE CONTENT"), http.StatusInternalServerError)
				return
			}

			contentType := http.DetectContentType(fileBytes)

			if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/gif" {
				helpers.JSONResponse(w, errors.New("UNSUPPORTED FILE FORMAT"), http.StatusBadRequest)
				return
			}

			os.MkdirAll("./views/assets/img/profile", os.ModePerm)

			fileName := fmt.Sprintf("%s-%s", helpers.UUID(), handler.Filename)
			picturePath := filepath.Join("./views/assets/img/profile", fileName)

			out, err := os.Create(picturePath)
			if err != nil {
				helpers.ErrorSimple(err)
				helpers.JSONResponse(w, errors.New("FAILED TO CREATE IMAGE FILE"), http.StatusInternalServerError)
				return
			}
			defer out.Close()

			_, err = out.Write(fileBytes)
			if err != nil {
				helpers.ErrorSimple(err)
				helpers.JSONResponse(w, errors.New("FAILED TO WRITE IMAGE CONTENT"), http.StatusInternalServerError)
				return
			}

			newUser.Values["picture"] = fileName

		}

		res := newUser.Create()

		// Connectez l'utilisateur après l'inscription
		controllers.CreateSession(w, r, fmt.Sprintf("%s", res.Values["id"]))
		controllers.BroadcastNewUser(strings.ReplaceAll(strings.Title(name), " ", "_"))

		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
	}
}
