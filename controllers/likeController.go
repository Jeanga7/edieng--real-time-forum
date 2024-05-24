package controllers

import (
	"errors"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func Like(w http.ResponseWriter, r *http.Request) {
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

	if len(entities) > 0 && r.FormValue("action") == "/like" {
		newLike := models.Like
		likedEntity := entities[0]
		AuthorLike := IdUser
		TypeLiked := r.FormValue("typeReacted")
		foundLiked := models.Like.Fetch([]interface{}{"authorId", "likedId", "typeLiked"}, []interface{}{AuthorLike, likedEntity.Values["id"], TypeLiked})
		foundDisliked := models.Dislike.Fetch([]interface{}{"authorId", "dislikedId", "typeDisliked"}, []interface{}{AuthorLike, likedEntity.Values["id"], TypeLiked})

		IsLogged, _ := Auth(r, w)
		if !IsLogged {
			return
		}

		if len(foundDisliked) != 0 {
			for _, dislike := range foundDisliked {
				dislike.Delete()
			}
		}

		if len(foundLiked) == 0 {
			newLike.Values = map[string]interface{}{
				"authorId":  AuthorLike,
				"likedId":   likedEntity.Values["id"],
				"typeLiked": TypeLiked,
			}
			newLike.Create()
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
			return
		} else {
			for _, like := range foundLiked {
				like.Delete()
			}

			// fmt.Printf("%s", foundID)
			// newLike.Delete("id", "=", fmt.Sprintf("%s", foundID))
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
			return
		}
	}
	helpers.JSONResponse(w, errors.New("NOTHING HAPPENED"), http.StatusAccepted)
}
