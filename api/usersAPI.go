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

var ListOfPossibleParams = []string{"firstname", "lastname", "nickname", "email", "password", "age", "gender"}

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch All Users
	var users []models.Model

	nbParams := 0
	for _, v := range ListOfPossibleParams {
		if len(r.FormValue(v)) > 0 {
			nbParams++
		}
	}
	paramName := ""
	expectedParams := true

	if nbParams == 1 {
		if len(r.FormValue("firstname")) > 0 {
			paramName = "firstname"
		} else if len(r.FormValue("lastname")) > 0 {
			paramName = "lastname"
		} else if len(r.FormValue("nickname")) > 0 {
			paramName = "nickname"
		} else if len(r.FormValue("email")) > 0 {
			paramName = "email"
		} else if len(r.FormValue("age")) > 0 {
			paramName = "age"
		} else if len(r.FormValue("gender")) > 0 {
			paramName = "gender"
		} else if len(r.FormValue("login")) > 0 {
			paramName = "login"
		} else {
			expectedParams = false
		}

		if expectedParams {
			if paramName == "login" {
				users = models.User.Fetch([]interface{}{"nickname"}, []interface{}{r.FormValue("login")})
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

	usersData := TransformModels(users, ListOfPossibleParams)
	WriteJSONResponse(w, usersData)
}

func ValidateFields(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)
	global.W = w

	form := html.EscapeString(r.FormValue("form"))
	if form == "login" {
		NicknameOrEmail := html.EscapeString(r.FormValue("login"))
		if len(NicknameOrEmail) > 0 {
			if strings.Contains(NicknameOrEmail, "@") && !helpers.ValidateInput("email", NicknameOrEmail) {
				helpers.JSONResponse(w, errors.New("INVALID EMAIL"), http.StatusBadRequest)
				return
			} else if !helpers.ValidateInput("nickname", NicknameOrEmail) {
				helpers.JSONResponse(w, errors.New("INVALID NICKNAME"), http.StatusBadRequest)
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

		user := models.User.Fetch([]interface{}{"email"}, []interface{}{NicknameOrEmail})
		if len(user) == 0 {
			user = models.User.Fetch([]interface{}{"nickname"}, []interface{}{NicknameOrEmail})
		}

		if len(user) == 1 {
			if len(password) > 0 {
				hashedPassword := fmt.Sprintf("%s", user[0].Values["password"])

				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
				if err == nil {
					controllers.CreateSession(w, r, fmt.Sprintf("%s", user[0].Values["id"]))
					controllers.BroadcastNewUser(user[0].Values["firstname"].(string) + " " + user[0].Values["lastname"].(string))
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
		validField := make(map[string]bool)
		validUser := true
		firstname := html.EscapeString(r.FormValue("firstname"))
		if len(firstname) > 0 {
			if !helpers.ValidateInput("name", firstname) {
				helpers.JSONResponse(w, errors.New("INVALID FIRSTNAME"), http.StatusBadRequest)
				return
			} else {
				validField["firstname"] = true
			}
		}

		lastname := html.EscapeString(r.FormValue("lastname"))
		if len(lastname) > 0 {
			if !helpers.ValidateInput("name", lastname) {
				helpers.JSONResponse(w, errors.New("INVALID LASTNAME"), http.StatusBadRequest)
				return
			} else {
				validField["lastname"] = true
			}
		}

		nickname := html.EscapeString(r.FormValue("nickname"))
		if len(nickname) > 0 {
			if !helpers.ValidateInput("nickname", nickname) {
				helpers.JSONResponse(w, errors.New("INVALID NICKNAME"), http.StatusBadRequest)
				return
			} else {
				validField["nickname"] = true
			}
		}

		email := html.EscapeString(r.FormValue("email"))
		if len(email) > 0 {
			if !helpers.ValidateInput("email", email) {
				helpers.JSONResponse(w, errors.New("INVALID EMAIL"), http.StatusBadRequest)
				return
			} else {
				validField["email"] = true
			}
		}

		password := html.EscapeString(r.FormValue("password"))
		if len(password) > 0 {
			if !helpers.ValidateInput("password", password) || len(password) < 8 || len(password) > 20 {
				helpers.JSONResponse(w, errors.New("INVALID PASSWORD"), http.StatusBadRequest)
				return
			} else {
				validField["password"] = true
			}
		}

		age := html.EscapeString(r.FormValue("age"))
		if len(age) > 0 {
			if !helpers.ValidateInput("age", age) {
				helpers.JSONResponse(w, errors.New("INVALID AGE"), http.StatusBadRequest)
				return
			} else {
				validField["age"] = true
			}
		}

		gender := html.EscapeString(r.FormValue("gender"))
		if len(gender) > 0 {
			if !helpers.ValidateInput("gender", gender) {
				helpers.JSONResponse(w, errors.New("INVALID GENDER"), http.StatusBadRequest)
				return
			} else {
				validField["gender"] = true
			}
		}

		// Vérifiez si l'email ou le nom d'utilisateur est déjà utilisé
		resEmail := models.User.Fetch([]interface{}{"email"}, []interface{}{email})
		if len(resEmail) > 0 {
			// Renvoyez une erreur si l'email ou le nom d'utilisateur est déjà utilisé
			helpers.JSONResponse(w, errors.New("EMAIL ALREADY USED"), http.StatusConflict)
			return
		}

		resNickname := models.User.Fetch([]interface{}{"nickname"}, []interface{}{nickname})
		if len(resNickname) > 0 {
			// Renvoyez une erreur si l'email ou le nom d'utilisateur est déjà utilisé
			helpers.JSONResponse(w, errors.New("NICKNAME ALREADY USED"), http.StatusConflict)
			return
		}

		for _, v := range ListOfPossibleParams {
			if !validField[v] {
				validUser = false
			}
		}

		if validUser {
			err := r.ParseMultipartForm(10 << 20)
			helpers.ErrorSimple(err)

			picture, handler, _ := r.FormFile("picture")

			// Si non, créez un nouvel utilisateur
			newUser := models.User
			newUser.Values = map[string]interface{}{
				"firstname": strings.ReplaceAll(strings.Title(firstname), " ", "_"),
				"lastname":  strings.ReplaceAll(strings.Title(lastname), " ", "_"),
				"nickname":  nickname,
				"email":     email,
				"age":       age,
				"gender":    gender,
				"password":  helpers.HashPassword(password),
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
			controllers.BroadcastNewUser(strings.ReplaceAll(strings.Title(firstname), " ", "_"))

			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		} else {
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		}
	}
}
