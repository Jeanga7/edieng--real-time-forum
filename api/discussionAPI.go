package api

import (
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"net/http"
)

func HandleDiscussions(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)
	global.W = w

	How := html.EscapeString(r.FormValue("type"))
	// Get the all the discussion between two users
	if How == "fetch" {

		with, user := html.EscapeString(r.FormValue("with")), html.EscapeString(r.FormValue("user"))

		WithFound, UserFound := models.User.Fetch([]interface{}{"nickname"}, []interface{}{with}), models.User.Fetch([]interface{}{"nickname"}, []interface{}{user})

		chatWith, chatUser := "", ""
		if len(WithFound) == 1 && len(UserFound) == 1 {
			chatWith, chatUser = (WithFound[0].Values["id"]).(string), (UserFound[0].Values["id"]).(string)
			discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatWith, chatUser})
			if len(discussions) == 0 {
				discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatUser, chatWith})
			}

			fields := append(models.Discussions.Columns, "updatedDate")
			usersData := TransformModels(discussions, fields)
			WriteJSONResponse(w, usersData)
		}
	}
}
