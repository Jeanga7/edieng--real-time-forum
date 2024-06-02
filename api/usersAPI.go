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
		} else {
			expectedParams = false
		}

		if expectedParams {
			users = models.User.Fetch([]interface{}{paramName}, []interface{}{params.Get(paramName)})
		} else {
			users = models.User.Fetch()
		}

	} else if len(params) == 2 && len(params.Get("username")) > 0 && len(params.Get("password")) > 0 {
		users = models.User.Fetch([]interface{}{"username"}, []interface{}{params.Get("username")})
		hashedPassword := fmt.Sprintf("%s", users[0].Values["password"])

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(params.Get("password")))
		if err != nil {
			users = []models.Model{}
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
