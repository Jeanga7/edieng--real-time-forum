package api

import (
	"forum/models"
	"net/http"
)

func FetchDislike(w http.ResponseWriter, r *http.Request) {
	var dislikes []models.Model
	params := r.URL.Query()

	if toShow := params.Get("dislikes"); len(toShow) > 0 && toShow == "all" {
		dislikes = models.Dislike.Fetch()
	}

	fields := []string{"id", "authorId", "dislikedId", "typeDisliked", "creationDate"}
	usersData := TransformModels(dislikes, fields)
	WriteJSONResponse(w, usersData)
}
