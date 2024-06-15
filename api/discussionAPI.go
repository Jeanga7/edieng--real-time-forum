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

	How := html.EscapeString(r.FormValue("type"))
	// Get the all the discussion between two users
	if How == "fetch" {

		chatUser := html.EscapeString(r.FormValue("with"))

		discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{chatUser, controllers.IdUser})
		if len(discussions) == 0 {
			discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{controllers.IdUser, chatUser})
		}

		fields := []string{"id", "idSender", "idReceiver", "messages", "creationDate"}
		usersData := TransformModels(discussions, fields)
		WriteJSONResponse(w, usersData)
		//if we are sending a message to a user
	} else if How == "send" {

		chatReceiver := html.EscapeString(r.FormValue("to"))

		Message := html.EscapeString(r.FormValue("message"))

		discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{controllers.IdUser, chatReceiver})

		if len(discussions) == 0 {
			controllers.SaveDiscussion(Message, controllers.IdUser, chatReceiver, helpers.Now())
		} else {
			discussionId := controllers.GetDiscussionId(controllers.IdUser, chatReceiver)
			controllers.AppendMessageToDiscussion(discussionId, Message, controllers.IdUser, chatReceiver, helpers.Now())
		}
	}
}

// be able to send message with curl
func SendMessage(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)
	global.W = w

}
