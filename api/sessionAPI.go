package api

import (
	"forum/models"
	"net/http"
)

func FetchSessions(w http.ResponseWriter, r *http.Request) {
	var sessions []models.Model
	params := r.URL.Query()

	if toShow := params.Get("session"); len(toShow) > 0 && toShow == "all" {
		sessions = models.Session.Fetch()
	}

	fields := []string{"id", "idUser", "expiryDate", "creationDate"}
	usersData := TransformModels(sessions, fields)
	WriteJSONResponse(w, usersData)
}
