package api

import (
	"forum/controllers"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"net/http"
)

func FetchDiscussions(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)
	global.W = w

	chatUser := html.EscapeString(r.FormValue("with"))

	discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatUser, controllers.IdUser})
	if len(discussions) == 0 {
		discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{controllers.IdUser, chatUser})
	}

	fields := []string{"id", "idSender", "idReceiver", "messages", "picture", "creationDate"}
	usersData := TransformModels(discussions, fields)
	WriteJSONResponse(w, usersData)
}
