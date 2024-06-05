package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	// var Posts []models.Model
	// params := r.URL.Query()

	Posts := models.Post.Fetch()

	usersData := func() []map[string]interface{} {
		allValues := []map[string]interface{}{}
		for _, u := range Posts {
			var data = make(map[string]interface{})
			data["id"] = u.Values["id"]
			data["authorId"] = u.Values["authorId"]
			data["title"] = u.Values["title"]
			data["content"] = u.Values["content"]
			data["categoriesNames"] = u.Values["categoriesNames"]
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
