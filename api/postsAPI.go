package api

import (
	"forum/database"
	"net/http"
)

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	posts := database.FetchDataJoin([]string{"posts", "users"}, []string{"*"}, "posts.authorId=users.id")

	for i := range posts {
		delete(posts[i], "password")
		delete(posts[i], "id")
		delete(posts[i], "authorId")
	}

	WriteJSONResponse(w, posts)
}
