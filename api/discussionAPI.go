package api

import (
	"errors"
	"forum/controllers"
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

		nickname := html.EscapeString(r.FormValue("with"))

		usersFound := models.User.Fetch([]interface{}{"nickname"}, []interface{}{nickname})

		chatUser := ""
		if len(usersFound) == 1 {
			chatUser = (usersFound[0].Values["id"]).(string)
			discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatUser, controllers.IdUser})
			if len(discussions) == 0 {
				discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{controllers.IdUser, chatUser})
			}

			usersData := TransformModels(discussions, models.Discussions.Columns)
			WriteJSONResponse(w, usersData)
		}
		//if we are sending a message to a user
	} else if How == "send" {
		nickname := html.EscapeString(r.FormValue("to"))
		usersFound := models.User.Fetch([]interface{}{"nickname"}, []interface{}{nickname})
		currentUser := models.User.Fetch([]interface{}{"id"}, []interface{}{controllers.IdUser})[0].Values["nickname"].(string)

		chatReceiver := ""
		if len(usersFound) == 1 {
			chatReceiver = (usersFound[0].Values["id"]).(string)
			Message := html.EscapeString(r.FormValue("message"))

			discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{controllers.IdUser, chatReceiver})
			if len(discussions) == 0 {
				discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatReceiver, controllers.IdUser})
			}

			if len(discussions) == 0 {
				controllers.SaveDiscussion(Message, controllers.IdUser, chatReceiver, helpers.Now())
			} else {
				discussionId := controllers.GetDiscussionId(controllers.IdUser, chatReceiver)
				controllers.AppendMessageToDiscussion(discussionId, Message, currentUser, nickname, helpers.Now())
			}
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		}
	}
}
