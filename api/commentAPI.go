package api

import (
	"forum/models"
	"net/http"
)

func FetchComments(w http.ResponseWriter, r *http.Request) {
	var comments []models.Model
	params := r.URL.Query()
	if toShow := params.Get("comments"); len(toShow) > 0 && toShow == "all" {
		comments = models.Comment.Fetch()
	}
	fields := []string{"id", "authorId", "postId", "content", "creationDate"}
	usersData := TransformModels(comments, fields)
	WriteJSONResponse(w, usersData)
}
