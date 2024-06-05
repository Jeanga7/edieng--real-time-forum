package api

import (
	"encoding/json"
	"fmt"
	"forum/helpers"
	"forum/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch All Users
	var users []models.Model

	listOfPossibleParams := []string{"id", "email", "name", "username", "login", "firstname", "lastname", "age", "gender"}

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

	usersData := func() []map[string]interface{} {
		allValues := []map[string]interface{}{}
		for _, u := range users {
			var data = make(map[string]interface{})
			data["id"] = u.Values["id"]
			data["name"] = u.Values["name"]
			data["username"] = u.Values["username"]
			data["email"] = u.Values["email"]
			data["picture"] = u.Values["picture"]
			data["creationDate"] = u.Values["creationDate"]
			allValues = append(allValues, data)
		}
		return allValues
	}()

	res, err := json.Marshal(usersData)
	helpers.ErrorSimple(err)

	// Header
	w.Header().Set("Content-Type", "application/json")

	// Response
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
