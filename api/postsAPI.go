package api

import (
	"forum/models"
	"net/http"
)

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Model
	params := r.URL.Query()

	if toShow := params.Get("posts"); len(toShow) > 0 && toShow == "all" {
		posts = models.Post.Fetch()
	}

	fields := []string{"id", "authorId", "title", "content", "categoriesNames", "picture", "creationDate"}
	usersData := TransformModels(posts, fields)
	WriteJSONResponse(w, usersData)
}
