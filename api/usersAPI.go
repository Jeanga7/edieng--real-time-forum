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

	params := r.URL.Query()
	paramName := ""
	expectedParams := true

	if len(params) == 1 {
		if len(params.Get("id")) > 0 {
			paramName = "id"
		} else if len(params.Get("email")) > 0 {
			paramName = "email"
		} else if len(params.Get("name")) > 0 {
			paramName = "name"
		} else if len(params.Get("username")) > 0 {
			paramName = "username"
		} else if len(params.Get("login")) > 0 {
			paramName = "login"
		} else {
			expectedParams = false
		}

		if expectedParams {
			if paramName == "login" {
				users = models.User.Fetch([]interface{}{"username"}, []interface{}{params.Get("login")})
				if len(users) == 0 {
					users = models.User.Fetch([]interface{}{"email"}, []interface{}{params.Get("login")})
				}
			} else {
				users = models.User.Fetch([]interface{}{paramName}, []interface{}{params.Get(paramName)})
			}
		} else {
			users = models.User.Fetch()
		}

	} else if len(params) == 2 && len(params.Get("login")) > 0 && len(params.Get("password")) > 0 {
		users = models.User.Fetch([]interface{}{"username"}, []interface{}{params.Get("login")})
		if len(users) == 0 {
			users = models.User.Fetch([]interface{}{"email"}, []interface{}{params.Get("login")})
		}

		if len(users) > 0 {
			hashedPassword := fmt.Sprintf("%s", users[0].Values["password"])

			err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(params.Get("password")))
			if err != nil {
				users = []models.Model{}
			}
		}
	}

	fields := []string{"id", "name", "username", "email", "picture", "creationDate"}
	usersData := TransformModels(users, fields)
	WriteJSONResponse(w, usersData)
}
