package controllers

import (
	"errors"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func Dislike(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)
	var entities []models.Model
	if r.FormValue("typeReacted") == "post" {
		entities = models.Post.Fetch([]interface{}{"id"}, []interface{}{r.FormValue("post")})
	} else {
		entities = models.Comment.Fetch([]interface{}{"id"}, []interface{}{r.FormValue("post")})
	}

	err := r.ParseForm()
	helpers.ErrorSimple(err)

	if len(entities) > 0 && r.FormValue("action") == "/dislike" {

		newDislike := models.Dislike
		dislikedEntity := entities[0]
		Authordislike := IdUser
		TypeDisliked := r.FormValue("typeReacted")
		foundDisliked := models.Dislike.Fetch([]interface{}{"authorId", "dislikedId", "typeDisliked"}, []interface{}{Authordislike, dislikedEntity.Values["id"], TypeDisliked})
		foundLiked := models.Like.Fetch([]interface{}{"authorId", "likedId", "typeLiked"}, []interface{}{Authordislike, dislikedEntity.Values["id"], TypeDisliked})

		IsLogged, _ := Auth(r, w)
		if !IsLogged {
			return
		}

		if len(foundLiked) != 0 {
			for _, like := range foundLiked {
				like.Delete()
			}
		}

		if len(foundDisliked) == 0 {
			newDislike.Values = map[string]interface{}{
				"authorId":     Authordislike,
				"dislikedId":   dislikedEntity.Values["id"],
				"typeDisliked": TypeDisliked,
			}
			newDislike.Create()
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
			return
		} else {
			for _, dislike := range foundDisliked {
				dislike.Delete()
			}

			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
			return
		}
	}
	helpers.JSONResponse(w, errors.New("NOTHING HAPPENED"), http.StatusAccepted)
}
