package api

import (
	"fmt"
	"forum/models"
	"net/http"

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

	} else if nbParams == 2 && len(r.FormValue("login")) > 0 && len(r.FormValue("password")) > 0 {
		users = models.User.Fetch([]interface{}{"username"}, []interface{}{r.FormValue("login")})
		if len(users) == 0 {
			users = models.User.Fetch([]interface{}{"email"}, []interface{}{r.FormValue("login")})
		}

		if len(users) > 0 {
			hashedPassword := fmt.Sprintf("%s", users[0].Values["password"])

			err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(r.FormValue("password")))
			if err != nil {
				users = []models.Model{}
			}
		}
	}

	fields := []string{"id", "name", "username", "email", "picture", "creationDate"}
	usersData := TransformModels(users, fields)
	WriteJSONResponse(w, usersData)
}
