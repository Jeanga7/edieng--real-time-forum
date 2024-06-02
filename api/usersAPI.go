package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func FetchUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch All Users
	users := models.User.Fetch()
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
