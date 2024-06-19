package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/controllers"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ListOfPossibleParams = []string{"firstname", "lastname", "nickname", "email", "password", "age", "gender"}

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch All Users
	var users []models.Model

	nbParams := 0
	cols := append(ListOfPossibleParams, "id")
	for _, v := range cols {
		if len(r.FormValue(v)) > 0 {
			nbParams++
		}
	}
	paramName := ""
	expectedParams := true

	if nbParams == 1 {
		if len(r.FormValue("id")) > 0 {
			paramName = "id"
		} else if len(r.FormValue("firstname")) > 0 {
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

	visibleCols := helpers.RemoveValueFromList(ListOfPossibleParams, "password")
	usersData := TransformModels(users, visibleCols)
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
					userInfos := user[0].Values
					userInfos["idSession"] = controllers.IdSession
					data, err := json.Marshal(userInfos)
					helpers.ErrorSimple(err)
					controllers.CreateSession(w, r, fmt.Sprintf("%s", user[0].Values["id"]))
					controllers.BroadcastNewUser(user[0].Values["firstname"].(string) + " " + user[0].Values["lastname"].(string))
					helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted, string(data))
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

		userData := map[string]map[string]interface{}{}

		for _, p := range ListOfPossibleParams {
			userData[p], userData[p] = map[string]interface{}{}, map[string]interface{}{}
		}

		userData["firstname"]["value"], userData["firstname"]["type"] = html.EscapeString(r.FormValue("firstname")), "name"
		userData["lastname"]["value"], userData["lastname"]["type"] = html.EscapeString(r.FormValue("lastname")), "name"
		userData["nickname"]["value"], userData["nickname"]["type"] = html.EscapeString(r.FormValue("nickname")), "nickname"
		userData["email"]["value"], userData["email"]["type"] = html.EscapeString(r.FormValue("email")), "email"
		userData["password"]["value"], userData["password"]["type"] = html.EscapeString(r.FormValue("password")), "password"
		userData["age"]["value"], userData["age"]["type"] = html.EscapeString(r.FormValue("age")), "age"
		userData["gender"]["value"], userData["gender"]["type"] = html.EscapeString(r.FormValue("gender")), "gender"

		for f, v := range userData {
			if !CheckField(w, v["value"].(string), f, v["type"].(string), validField) {
				return
			}
		}

		// Vérifiez si l'email ou le nom d'utilisateur est déjà utilisé
		resEmail := models.User.Fetch([]interface{}{"email"}, []interface{}{userData["email"]["value"]})
		if len(resEmail) > 0 {
			// Renvoyez une erreur si l'email ou le nom d'utilisateur est déjà utilisé
			helpers.JSONResponse(w, errors.New("EMAIL ALREADY USED"), http.StatusConflict)
			return
		}

		resNickname := models.User.Fetch([]interface{}{"nickname"}, []interface{}{userData["nickname"]["value"]})
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
			newUser := models.User
			newUser.Values = map[string]interface{}{
				"firstname": userData["firstname"]["value"].(string),
				"lastname":  userData["lastname"]["value"].(string),
				"nickname":  userData["nickname"]["value"].(string),
				"email":     userData["email"]["value"].(string),
				"age":       userData["age"]["value"],
				"gender":    userData["gender"]["value"].(string),
				"password":  helpers.HashPassword(userData["password"]["value"].(string)),
			}

			res := newUser.Create()

			// Connectez l'utilisateur après l'inscription
			userInfos := res.Values
			userInfos["idSession"] = controllers.IdSession
			data, err := json.Marshal(userInfos)
			helpers.ErrorSimple(err)
			controllers.CreateSession(w, r, fmt.Sprintf("%s", res.Values["id"]))
			controllers.BroadcastNewUser(userData["firstname"]["value"].(string))
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted, string(data))
		} else {
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		}
	}
}

func CheckField(w http.ResponseWriter, fieldVar, fieldName, fieldType string, validField map[string]bool) bool {
	if len(fieldVar) > 0 {
		if !helpers.ValidateInput(fieldType, fieldVar) {
			helpers.JSONResponse(w, errors.New("INVALID "+strings.ToUpper(fieldName)), http.StatusBadRequest)
			return false
		} else {
			validField[fieldName] = true
		}
	}
	return true
}

// =======================AJOUT=====================================

func GetAllSortedUsers(w http.ResponseWriter, r *http.Request) {
	Usernames, Sessions, shouldReturn, _ := controllers.FindUsersOnline(r)
	if shouldReturn {
		return
	}
	All_users, _ := controllers.GetUsers(Usernames)
	SortedUsers := controllers.SortUsers(All_users, Sessions, Usernames)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SortedUsers)
}
