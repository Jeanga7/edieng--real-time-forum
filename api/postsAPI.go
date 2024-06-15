package api

import (
	"forum/controllers"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"net/http"
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
