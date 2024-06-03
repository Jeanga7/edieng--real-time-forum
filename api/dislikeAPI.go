package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func FetchDislike(w http.ResponseWriter, r *http.Request) {
	var likes []models.Model
	params := r.URL.Query()

	if toShow := params.Get("dislikes"); len(toShow) > 0 && toShow == "all" {
		likes = models.Dislike.Fetch()
	}

	usersData := func() []map[string]interface{} {
		allValues := []map[string]interface{}{}
		for _, u := range likes {
			var data = make(map[string]interface{})
			data["id"] = u.Values["id"]
			data["authorId"] = u.Values["authorId"]
			data["dislikedId"] = u.Values["dislikedId"]
			data["typeDisliked"] = u.Values["typeDisliked"]
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
