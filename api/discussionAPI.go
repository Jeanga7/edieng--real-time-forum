package api

import (
	"forum/models"
	"net/http"
)

func FetchDiscussions(w http.ResponseWriter, r *http.Request) {
	var discussions []models.Model
	params := r.URL.Query()

	if toShow := params.Get("discussions"); len(toShow) > 0 && toShow == "all" {
		discussions = models.Discussions.Fetch()
	}

	fields := []string{"id", "idSender", "idReceiver", "messages", "picture", "creationDate"}
	usersData := TransformModels(discussions, fields)
	WriteJSONResponse(w, usersData)
}
