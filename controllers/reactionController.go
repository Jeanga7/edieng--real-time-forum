package controllers

import (
	"errors"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func React(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)

	IsLogged, _ := Auth(r, w)
	if !IsLogged {
		helpers.JSONResponse(w, errors.New("SESSION NOT FOUND"), http.StatusBadRequest)
		return
	}

	action, reaction, AuthorReaction, TypeReacted, entityId := r.FormValue("action"), r.FormValue("reaction"), r.FormValue("author"), r.FormValue("typeReacted"), r.FormValue("post")
	like, dislike := (reaction == "like"), (reaction == "dislike")
	emptyAction, emptyReaction, emptyAuthor, emptyType := EmptyStr(action), EmptyStr(reaction), EmptyStr(AuthorReaction), EmptyStr(TypeReacted)

	if (!like && !dislike) || emptyAction || emptyReaction || emptyAuthor || emptyType {
		helpers.JSONResponse(w, errors.New("SOME FIELDS ARE EMPTY"), http.StatusBadRequest)
		return
	}

	var entities []models.Model
	if TypeReacted == "post" {
		entities = models.Post.Fetch([]interface{}{"id"}, []interface{}{entityId})
	} else if TypeReacted == "comment" {
		entities = models.Comment.Fetch([]interface{}{"id"}, []interface{}{entityId})
	}

	users := models.User.Fetch([]interface{}{"nickname"}, []interface{}{AuthorReaction})

	objectNotFound, userNotFound := len(entities) == 0, len(users) == 0

	if objectNotFound {
		helpers.JSONResponse(w, errors.New("POST OR COMMENT NOT FOUND"), http.StatusBadRequest)
		return
	} else if userNotFound {
		helpers.JSONResponse(w, errors.New("USER NOT FOUND"), http.StatusBadRequest)
		return
	}

	userId := users[0].Values["id"]
	reactedEntity := entities[0]
	err := r.ParseForm()
	helpers.ErrorSimple(err)

	var reactIdString, reactTypeString, revReactIdString, revReactTypeString string
	var foundReaction, foundRevReaction []models.Model
	var newReaction models.Model

	if like {
		newReaction = models.Like
		reactIdString, reactTypeString = "likedId", "typeLiked"
		revReactIdString, revReactTypeString = "dislikedId", "typeDisliked"
		foundReaction = models.Like.Fetch([]interface{}{"authorId", reactIdString, reactTypeString}, []interface{}{userId, reactedEntity.Values["id"], TypeReacted})
		foundRevReaction = models.Dislike.Fetch([]interface{}{"authorId", revReactIdString, revReactTypeString}, []interface{}{userId, reactedEntity.Values["id"], TypeReacted})
	} else {
		newReaction = models.Dislike
		reactIdString, reactTypeString = "dislikedId", "typeDisliked"
		revReactIdString, revReactTypeString = "likedId", "typeLiked"
		foundReaction = models.Dislike.Fetch([]interface{}{"authorId", reactIdString, reactTypeString}, []interface{}{userId, reactedEntity.Values["id"], TypeReacted})
		foundRevReaction = models.Like.Fetch([]interface{}{"authorId", revReactIdString, revReactTypeString}, []interface{}{userId, reactedEntity.Values["id"], TypeReacted})
	}

	if len(foundRevReaction) != 0 {
		for _, r := range foundRevReaction {
			r.Delete()
		}
	}

	if len(foundReaction) == 0 {
		newReaction.Values = map[string]interface{}{
			"authorId":      userId,
			reactIdString:   reactedEntity.Values["id"],
			reactTypeString: TypeReacted,
		}
		newReaction.Create()
		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		return
	} else {
		for _, r := range foundReaction {
			r.Delete()
		}

		// fmt.Printf("%s", foundID)
		// newReaction.Delete("id", "=", fmt.Sprintf("%s", foundID))
		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
		return
	}
}

func EmptyStr(el string) bool {
	return len(el) == 0
}
