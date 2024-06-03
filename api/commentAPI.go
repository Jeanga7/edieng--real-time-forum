package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func FetchComments(w http.ResponseWriter, r *http.Request) {
	var comments []models.Model
	params := r.URL.Query()

	if toShow := params.Get("comments"); len(toShow) > 0 && toShow == "all" {
		comments = models.Comment.Fetch()
	}

	usersData := func() []map[string]interface{} {
		allValues := []map[string]interface{}{}
		for _, u := range comments {
			var data = make(map[string]interface{})
			data["id"] = u.Values["id"]
			data["authorId"] = u.Values["authorId"]
			data["postId"] = u.Values["postId"]
			data["content"] = u.Values["content"]
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
