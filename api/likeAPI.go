package api

import (
	"forum/models"
	"net/http"
)

func FetchLike(w http.ResponseWriter, r *http.Request) {
	var likes []models.Model
	params := r.URL.Query()

	if toShow := params.Get("likes"); len(toShow) > 0 && toShow == "all" {
		likes = models.Like.Fetch()
	}

	fields := []string{"id", "authorId", "likedId", "typeLiked", "creationDate"}
	usersData := TransformModels(likes, fields)
	WriteJSONResponse(w, usersData)
}
