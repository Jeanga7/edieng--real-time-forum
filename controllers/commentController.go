package controllers

import (
	"errors"
	"fmt"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"net/http"
)

var (
	NumCommentLike    int
	NumCommentDisLike int
)

// Fonction pour obtenir les commentaires d'un poste
func GetPostComments(comments []models.Model) []helpers.PostComment {
	var postComments []helpers.PostComment
	NumCommentLike = 0

	for _, comment := range comments {
		var newComment helpers.PostComment
		newComment.Id = helpers.InterfaceToString(comment.Values["id"])
		newComment.Content = helpers.InterfaceToString(comment.Values["content"])

		likes := models.Like.Fetch([]interface{}{"likedId", "typeLiked"}, []interface{}{comment.Values["id"], "comment"})
		newComment.Likes = len(likes)
		NumCommentLike += len(likes)
		if len(models.Like.Fetch([]interface{}{"likedId", "typeLiked", "authorId"}, []interface{}{comment.Values["id"], "comment", IdUser})) > 0 {
			newComment.Liked = true
		}

		dislikes := models.Dislike.Fetch([]interface{}{"dislikedId", "typeDisliked"}, []interface{}{comment.Values["id"], "comment"})
		newComment.Dislikes = len(dislikes)
		NumCommentDisLike = len(dislikes)
		if len(models.Dislike.Fetch([]interface{}{"dislikedId", "typeDisliked", "authorId"}, []interface{}{comment.Values["id"], "comment", IdUser})) > 0 {
			newComment.Disliked = true
		}

		authorID := comment.Values["authorId"]
		authors := models.User.Fetch([]interface{}{"id"}, []interface{}{authorID})
		if len(authors) > 0 {
			author := authors[0]
			newComment.AuthorName = helpers.InterfaceToString(author.Values["name"])
			newComment.AuthorUsername = helpers.InterfaceToString(author.Values["username"])
			newComment.AuthorPicture = helpers.DefaultIfNil(author.Values["picture"])
		}

		postComments = append(postComments, newComment)
	}
	return postComments
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)
	IsLogged, _ := Auth(r, w)
	if !IsLogged {
		return
	}

	postId := r.FormValue("postId")
	posts := models.Post.Fetch([]interface{}{"id"}, []interface{}{postId})

	if len(posts) == 0 {
		helpers.JSONResponse(w, errors.New("POST NOT FOUND"), http.StatusBadRequest)
		return
	}

	newComment := models.Comment

	content := html.EscapeString(r.FormValue("commentTxt"))
	if limit := 500; len(content) > limit {
		helpers.JSONResponse(w, fmt.Errorf("COMMENT EXCEEDS %d CHARS", limit), http.StatusBadRequest)
		return
	}

	// if !helpers.ValidateInput("content", content) {
	// 	helpers.JSONResponse(w, errors.New("INVALID CONTENT"), http.StatusBadRequest)
	// 	return
	// }

	authorComment := IdUser

	newComment.Values = map[string]interface{}{
		"authorId": authorComment,
		"postId":   postId,
		"content":  content,
	}

	newComment.Create()

	helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
}

func CheckLogging(author interface{}, r *http.Request, w http.ResponseWriter) bool {
	global.W = w
	session := models.Session.Fetch([]interface{}{"idUser"}, []interface{}{author})
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return true
	}
	if len(session) == 0 || sessionCookie == nil {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		helpers.ErrorHandler(w, err, "Session not found", http.StatusUnauthorized)
		return true
	}
	return false
}
