package api

import (
	"errors"
	"fmt"
	"forum/controllers"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"net/http"
	"strings"
)

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	posts := ReadPost(w, r)
	WriteJSONResponse(w, posts)
}

func ReadPost(w http.ResponseWriter, r *http.Request) []helpers.ForumPost {
	global.W = w
	controllers.LikedPost = []helpers.ForumPost{}
	controllers.DisLikedPost = []helpers.ForumPost{}
	controllers.MyPosts = []helpers.ForumPost{}
	controllers.MyComments = []helpers.ForumPost{}
	controllers.CountComment = 0
	controllers.NumCommentLike = 0
	controllers.NumCommentDisLike = 0

	helpers.AllowedMethod(w, "GET", r)
	var posts []models.Model

	if params := r.URL.Query(); len(params.Get("category")) > 0 {
		posts = models.Post.FetchLike([]interface{}{"categoriesNames"}, []interface{}{"%" + params.Get("category") + "%"})
	} else if len(params.Get("search")) > 0 {
		posts = controllers.FillPostSearched(params, posts)
	} else {
		posts = models.Post.Fetch()
	}

	var res []helpers.ForumPost
	if len(posts) > 0 {
		for _, post := range posts {
			authors := models.User.Fetch([]interface{}{"id"}, []interface{}{post.Values["authorId"]})
			if len(authors) > 0 {
				author := authors[0]

				var newValues helpers.ForumPost
				newValues.Id = helpers.InterfaceToString(post.Values["id"])
				newValues.AuthorId = helpers.InterfaceToString(post.Values["authorId"])
				newValues.AuthorName = helpers.InterfaceToString(author.Values["firstname"].(string) + " " + author.Values["lastname"].(string))
				newValues.AuthorNickname = helpers.InterfaceToString(author.Values["nickname"])
				newValues.Title = helpers.InterfaceToString(post.Values["title"])
				newValues.Content = helpers.InterfaceToString(post.Values["content"])
				newValues.Categories = helpers.InterfaceToString(post.Values["categoriesNames"])
				newValues.PostPicture = helpers.EmptyStringIfNil(post.Values["picture"])
				newValues.AuthorPicture = helpers.DefaultIfNil(author.Values["picture"])
				newValues.CreationDate = helpers.InterfaceToString(post.Values["creationDate"])

				// Récupération des commentaires du post
				comments := models.Comment.Fetch([]interface{}{"postId"}, []interface{}{newValues.Id})
				newValues.Comments = controllers.GetPostComments(comments)

				// Récupération des likes du post
				likes := models.Like.Fetch([]interface{}{"likedId", "typeLiked"}, []interface{}{newValues.Id, "post"})
				newValues.Likes = len(likes)

				// Récupération des dislikes du post
				dislikes := models.Dislike.Fetch([]interface{}{"dislikedId", "typeDisliked"}, []interface{}{newValues.Id, "post"})
				newValues.Dislikes = len(dislikes)

				// Récupération du nombre total de commentaires du post
				newValues.NumComments = len(comments)

				// Check if liked
				OwnerLike := models.Like.Fetch([]interface{}{"authorId", "likedId"}, []interface{}{controllers.IdUser, post.Values["id"]})
				if len(OwnerLike) > 0 {
					newValues.Liked = true
					controllers.LikedPost = append(controllers.LikedPost, newValues)
				}

				// Check if disliked
				OwnerDislike := models.Dislike.Fetch([]interface{}{"authorId", "dislikedId"}, []interface{}{controllers.IdUser, post.Values["id"]})
				if len(OwnerDislike) > 0 {
					newValues.Disliked = true
					controllers.DisLikedPost = append(controllers.DisLikedPost, newValues)
				}

				//Get the User's posts
				OwnerPost := models.Post.Fetch([]interface{}{"authorId"}, []interface{}{controllers.IdUser})
				if len(OwnerPost) > 0 {
					if newValues.AuthorId == OwnerPost[0].Values["authorId"] {
						controllers.MyPosts = append(controllers.MyPosts, newValues)
					}
				}
				//Get the User's comments
				OwnerComment := models.Comment.Fetch([]interface{}{"authorId", "postId"}, []interface{}{controllers.IdUser, newValues.Id})
				if len(OwnerComment) > 0 {
					controllers.MyComments = append(controllers.MyComments, newValues)
					controllers.CountComment += len(OwnerComment)
				}
				res = append(res, newValues)
			}
		}
	}
	return res
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)

	author := html.EscapeString(r.FormValue("author"))
	if len(author) == 0 {
		helpers.JSONResponse(w, errors.New("MISSING POST AUTHOR"), http.StatusBadRequest)
		return
	}

	sessions := models.Session.Fetch([]interface{}{"id"}, []interface{}{author})
	if len(sessions) == 0 {
		helpers.JSONResponse(w, errors.New("WRONG POST AUTHOR"), http.StatusBadRequest)
		return
	}

	authorID := sessions[0].Values["idUser"].(string)

	title := html.EscapeString(r.FormValue("title"))
	if limit := 100; len(title) > limit {
		helpers.JSONResponse(w, fmt.Errorf("TITLE EXCEEDS %d CHARS", limit), http.StatusBadRequest)
		return
	}

	if len(title) > 0 {
		if !helpers.ValidateInput("title", title) {
			helpers.JSONResponse(w, errors.New("POST TITLE ERROR"), http.StatusBadRequest)
			return
		}
	}

	content := html.EscapeString(r.FormValue("content"))
	if limit := 1500; len(content) > limit {
		helpers.JSONResponse(w, fmt.Errorf("TEXT EXCEEDS %d CHARS", limit), http.StatusBadRequest)
		return
	}

	category := html.EscapeString(r.FormValue("category"))
	if len(category) > 0 {
		foundCategory := models.Category.Fetch([]interface{}{"name"}, []interface{}{category})
		if len(foundCategory) == 0 {
			if !helpers.ValidateInput("category", category) {
				helpers.JSONResponse(w, errors.New("INVALID NEW CATEGORY"), http.StatusBadRequest)
				return
			}
		}
	}

	allFilled := notEmptyStr(authorID) && notEmptyStr(title) && notEmptyStr(content) && notEmptyStr(category)

	// Create Post
	if allFilled {
		newCategory := models.Category
		newCategory.Values = make(map[string]interface{})
		newCategory.Values["name"] = html.EscapeString(category)
		newCategory.Create()

		newPost := models.Post
		newPost.Values = map[string]interface{}{
			"authorId":        authorID,
			"title":           title,
			"content":         content,
			"categoriesNames": strings.Title(category),
		}
		newPost.Create()
	}

	helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
}

func notEmptyStr(str string) bool {
	return str != "" && str != "undefined"
}
