package api

import (
	"forum/models"
	"net/http"
)

func FetchCategory(w http.ResponseWriter, r *http.Request) {
	var posts []models.Model
	params := r.URL.Query()

	if toShow := params.Get("category"); len(toShow) > 0 {
		posts = models.Post.FetchLike([]interface{}{"categoriesNames"}, []interface{}{"%" + params.Get("category") + "%"})
	}

	fields := []string{"id", "authorId", "title", "content", "categoriesNames", "picture", "creationDate"}
	usersData := TransformModels(posts, fields)

	WriteJSONResponse(w, usersData)
}
