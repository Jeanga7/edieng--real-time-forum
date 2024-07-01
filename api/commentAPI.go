package api

import (
	"errors"
	"fmt"
	"forum/controllers"
	"forum/database"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"net/http"
)

func FetchComments(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)
	global.W = w

	postId := r.FormValue("postId")
	var res []map[string]interface{}

	if len(postId) > 0 {
		res = database.FetchDataJoin(
			[]string{"comments", "users"},
			[]string{"*", "comments.creationDate"},
			"comments.authorId=users.id AND comments.postId = ?",
			postId,
		)
	}

	WriteJSONResponse(w, res)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)
	IsLogged, _ := controllers.Auth(r, w)
	if !IsLogged {
		return
	}

	postId, sessionId, commentTxt := r.FormValue("postId"), r.FormValue("sessionId"), r.FormValue("commentTxt")

	if len(postId) > 0 && len(sessionId) > 0 && len(commentTxt) > 0 {
		posts := models.Post.Fetch([]interface{}{"id"}, []interface{}{postId})
		if len(posts) == 0 {
			helpers.JSONResponse(w, errors.New("POST NOT FOUND"), http.StatusBadRequest)
			return
		}

		sessions := models.Session.Fetch([]interface{}{"id"}, []interface{}{sessionId})
		if len(sessions) == 0 {
			helpers.JSONResponse(w, errors.New("SESSION NOT FOUND"), http.StatusBadRequest)
			return
		}

		newComment := models.Comment

		content := html.EscapeString(r.FormValue("commentTxt"))
		if limit := 500; len(content) > limit {
			helpers.JSONResponse(w, fmt.Errorf("COMMENT EXCEEDS %d CHARS", limit), http.StatusBadRequest)
			return
		}

		authorComment := sessions[0].Values["idUser"].(string)

		newComment.Values = map[string]interface{}{
			"authorId": authorComment,
			"postId":   postId,
			"content":  content,
		}

		newComment.Create()

		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		return
	}
	helpers.JSONResponse(w, errors.New("MISSING DATA"), http.StatusAccepted)
}
