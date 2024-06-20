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

			usersData := TransformModels(discussions, models.Discussions.Columns)
			WriteJSONResponse(w, usersData)
		}
		//if we are sending a message to a user
	}
	// else if How == "send" {
	// 	nickname := html.EscapeString(r.FormValue("to"))
	// 	usersFound := models.User.Fetch([]interface{}{"nickname"}, []interface{}{nickname})
	// 	currentUser := models.User.Fetch([]interface{}{"id"}, []interface{}{controllers.IdUser})[0].Values["nickname"].(string)

	// 	chatReceiver := ""
	// 	if len(usersFound) == 1 {
	// 		chatReceiver = (usersFound[0].Values["id"]).(string)
	// 		message := html.EscapeString(r.FormValue("message"))

	// 		discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{controllers.IdUser, chatReceiver})
	// 		if len(discussions) == 0 {
	// 			discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatReceiver, controllers.IdUser})
	// 		}

	// 		if len(discussions) == 0 {
	// 			controllers.SaveDiscussion(message, controllers.IdUser, chatReceiver, helpers.Now())
	// 		} else {
	// 			discussionId := controllers.GetDiscussionId(controllers.IdUser, chatReceiver)
	// 			controllers.AppendMessageToDiscussion(discussionId, message, currentUser, nickname, helpers.Now())
	// 		}

	// 		err := controllers.StoreMessage(controllers.IdUser, currentUser, chatReceiver, nickname, message)

	// 		if err != nil {
	// 			helpers.JSONResponse(w, err, http.StatusInternalServerError)
	// 		}
	// 		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
	// 	}
	// }
}
