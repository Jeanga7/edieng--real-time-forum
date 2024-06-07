package controllers

import (
	"errors"
	"fmt"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	LikedPost    []helpers.ForumPost
	DisLikedPost []helpers.ForumPost
	MyPosts      []helpers.ForumPost
	MyComments   []helpers.ForumPost
	CountComment int
)

func ReadPost(w http.ResponseWriter, r *http.Request) []helpers.ForumPost {
	global.W = w
	LikedPost = []helpers.ForumPost{}
	DisLikedPost = []helpers.ForumPost{}
	MyPosts = []helpers.ForumPost{}
	MyComments = []helpers.ForumPost{}
	CountComment = 0
	NumCommentLike = 0
	NumCommentDisLike = 0

	helpers.AllowedMethod(w, "GET", r)
	var posts []models.Model

	if params := r.URL.Query(); len(params.Get("category")) > 0 {
		posts = models.Post.FetchLike([]interface{}{"categoriesNames"}, []interface{}{"%" + params.Get("category") + "%"})
	} else if len(params.Get("search")) > 0 {
		posts = FillPostSearched(params, posts)
	} else {
		// Vérification pour les routes inconnues
		if r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/?") {
			global.Status = 404
			return nil
		}
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
				newValues.Comments = GetPostComments(comments)

				// Récupération des likes du post
				likes := models.Like.Fetch([]interface{}{"likedId", "typeLiked"}, []interface{}{newValues.Id, "post"})
				newValues.Likes = len(likes)

				// Récupération des dislikes du post
				dislikes := models.Dislike.Fetch([]interface{}{"dislikedId", "typeDisliked"}, []interface{}{newValues.Id, "post"})
				newValues.Dislikes = len(dislikes)

				// Récupération du nombre total de commentaires du post
				newValues.NumComments = len(comments)

				// Check if liked
				OwnerLike := models.Like.Fetch([]interface{}{"authorId", "likedId"}, []interface{}{IdUser, post.Values["id"]})
				if len(OwnerLike) > 0 {
					newValues.Liked = true
					LikedPost = append(LikedPost, newValues)
				}

				// Check if disliked
				OwnerDislike := models.Dislike.Fetch([]interface{}{"authorId", "dislikedId"}, []interface{}{IdUser, post.Values["id"]})
				if len(OwnerDislike) > 0 {
					newValues.Disliked = true
					DisLikedPost = append(DisLikedPost, newValues)
				}

				//Get the User's posts
				OwnerPost := models.Post.Fetch([]interface{}{"authorId"}, []interface{}{IdUser})
				if len(OwnerPost) > 0 {
					if newValues.AuthorId == OwnerPost[0].Values["authorId"] {
						MyPosts = append(MyPosts, newValues)
					}
				}
				//Get the User's comments
				OwnerComment := models.Comment.Fetch([]interface{}{"authorId", "postId"}, []interface{}{IdUser, newValues.Id})
				if len(OwnerComment) > 0 {
					MyComments = append(MyComments, newValues)
					CountComment += len(OwnerComment)
				}
				res = append(res, newValues)
			}
		}
	}
	return res
}

func FillPostSearched(params url.Values, posts []models.Model) []models.Model {
	wanted := params.Get("search")
	wantedCategories := models.Post.FetchLike([]interface{}{"categoriesNames"}, []interface{}{"%" + wanted + "%"})
	wantedTitle := models.Post.FetchLike([]interface{}{"title"}, []interface{}{"%" + wanted + "%"})
	wantedContent := models.Post.FetchLike([]interface{}{"content"}, []interface{}{"%" + wanted + "%"})

	Name := models.User.FetchLike([]interface{}{"firstname"}, []interface{}{"%" + wanted + "%"})
	if len(Name) == 0 {
		Name = models.User.FetchLike([]interface{}{"lastname"}, []interface{}{"%" + wanted + "%"})
	}
	Username := models.User.FetchLike([]interface{}{"nickname"}, []interface{}{"%" + wanted + "%"})

	var tabWantedNames [][]models.Model
	for _, v := range Name {
		tabWantedNames = append(tabWantedNames, models.Post.FetchLike([]interface{}{"authorId"}, []interface{}{"%" + v.Values["id"].(string) + "%"}))
	}
	var tabWantedNicknames [][]models.Model
	for _, v := range Username {
		tabWantedNicknames = append(tabWantedNicknames, models.Post.FetchLike([]interface{}{"authorId"}, []interface{}{"%" + v.Values["id"].(string) + "%"}))
	}
	if len(wantedCategories) > 0 {
		posts = append(posts, wantedCategories...)
	}
	if len(wantedTitle) > 0 {
		posts = append(posts, wantedTitle...)
	}

	if len(wantedContent) > 0 {
		posts = append(posts, wantedContent...)
	}

	for _, tabWantedName := range tabWantedNames {
		if len(tabWantedName) > 0 {
			posts = append(posts, tabWantedName...)
		}
	}

	for _, tabWantedNickname := range tabWantedNicknames {
		if len(tabWantedNickname) > 0 {
			posts = append(posts, tabWantedNickname...)
		}
	}

	var finalPosts []models.Model
	var MatchedPost = make(map[string]bool)
	for _, post := range posts {
		if !MatchedPost[post.Values["id"].(string)] {
			finalPosts = append(finalPosts, post)
		}
		MatchedPost[post.Values["id"].(string)] = true
	}

	return finalPosts
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)

	title := html.EscapeString(r.FormValue("title"))
	if limit := 100; len(title) > limit {
		helpers.JSONResponse(w, fmt.Errorf("TITLE EXCEEDS %d CHARS", limit), http.StatusBadRequest)
		return
	}

	if !helpers.ValidateInput("title", title) {
		helpers.JSONResponse(w, errors.New("POST TITLE ERROR"), http.StatusBadRequest)
		return
	}

	// r.ParseMultipartForm(10 << 20)

	content := html.EscapeString(r.FormValue("content"))
	if limit := 1500; len(content) > limit {
		helpers.JSONResponse(w, fmt.Errorf("TEXT EXCEEDS %d CHARS", limit), http.StatusBadRequest)
		return
	}

	picture, handler, err := r.FormFile("picture")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println(picture)
		fmt.Println(handler)
		fmt.Println(err)
		helpers.JSONResponse(w, errors.New("FILE FORMAT ERROR"), http.StatusBadRequest)
		return
	}

	if picture != nil {
		defer picture.Close()
	}

	const maxSize = 20 << 20 // 20 MB
	if handler != nil && handler.Size > maxSize {
		helpers.JSONResponse(w, errors.New("FILE EXCEEDS 20MB"), http.StatusBadRequest)
		return
	}

	if content == "" && picture == nil {
		helpers.JSONResponse(w, errors.New("MISSING CONTENT OR PICTURE"), http.StatusBadRequest)
		return
	}

	IsLogged, _ := Auth(r, w)
	if !IsLogged {
		return
	}

	// Check new categories
	for _, cat := range r.Form["categories[]"] {
		found := models.Category.Fetch([]interface{}{"name"}, []interface{}{cat})
		if len(found) == 0 {
			if !helpers.ValidateInput("category", cat) {
				helpers.ErrorHandler(w, errors.New("INVALID NEW CATEGORY"), "ValidateInput", 400)
				return
			}

			newCategory := models.Category
			newCategory.Values = make(map[string]interface{})
			newCategory.Values["name"] = html.EscapeString(cat)
			newCategory.Create()
		}
	}

	// Accédez aux valeurs cochées
	categories := strings.Join(r.Form["categories[]"], ",")
	if len(categories) == 0 {
		helpers.JSONResponse(w, errors.New("MISSING POST CATEGORIES"), http.StatusBadRequest)
		return
	}

	newPost := models.Post
	newPost.Values = map[string]interface{}{
		"authorId":        IdUser,
		"title":           title,
		"content":         content,
		"categoriesNames": categories,
	}

	if picture != nil {
		fileBytes, err := io.ReadAll(picture)
		if err != nil {
			helpers.ErrorSimple(err)
			helpers.JSONResponse(w, errors.New("FAILED TO READ IMAGE CONTENT"), http.StatusInternalServerError)
			return
		}

		contentType := http.DetectContentType(fileBytes)
		if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/gif" && !strings.Contains(contentType, "text/xml") && contentType != "image/svg+xml" {
			helpers.JSONResponse(w, errors.New("UNSUPPORTED FILE FORMAT"), http.StatusBadRequest)
			return
		}

		os.MkdirAll(helpers.Storage, os.ModePerm)

		fileName := fmt.Sprintf("%s-%s", helpers.UUID(), handler.Filename)
		picturePath := filepath.Join(helpers.Storage, fileName)

		out, err := os.Create(picturePath)
		if err != nil {
			helpers.ErrorSimple(err)
			helpers.JSONResponse(w, errors.New("FAILED TO CREATE IMAGE FILE"), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = out.Write(fileBytes)
		if err != nil {
			helpers.ErrorSimple(err)
			helpers.JSONResponse(w, errors.New("FAILED TO WRITE IMAGE CONTENT"), http.StatusInternalServerError)
			return
		}

		newPost.Values["picture"] = fileName
	}

	newPost.Create()
	helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)
	posts := models.Post.Fetch()
	postID := posts[helpers.RandN(0, 1)]
	postToEdit := models.Post.Fetch([]interface{}{"id"}, []interface{}{postID.Values["id"]})

	if len(postToEdit) != 0 {
		title := html.EscapeString(r.FormValue("title"))
		if !helpers.ValidateInput("title", title) {
			helpers.JSONResponse(w, errors.New("POST TITLE ERROR"), http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			helpers.JSONResponse(w, errors.New("FILE SIZE ERROR"), http.StatusBadRequest)
			return
		}

		content := html.EscapeString(r.FormValue("content"))
		picture, handler, _ := r.FormFile("picture")

		if content == "" && picture == nil {
			helpers.JSONResponse(w, errors.New("MISSING CONTENT OR PICTURE"), http.StatusBadRequest)
			return
		}

		if postToEdit == nil {
			helpers.JSONResponse(w, errors.New("POST TITLE ERROR"), http.StatusBadRequest)
			return
		}

		for _, champ := range postToEdit {
			dataToUpdate := map[string]interface{}{
				"title":   title,
				"content": content,
			}

			if picture != nil {
				fileBytes, err := io.ReadAll(picture)
				if err != nil {
					helpers.ErrorSimple(err)
					helpers.JSONResponse(w, errors.New("FAILED TO READ IMAGE CONTENT"), http.StatusInternalServerError)
					return
				}

				contentType := http.DetectContentType(fileBytes)

				if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/gif" {
					helpers.JSONResponse(w, errors.New("UNSUPPORTED FILE FORMAT"), http.StatusBadRequest)
					return
				}

				os.MkdirAll(helpers.Storage, os.ModePerm)

				fileName := fmt.Sprintf("%s-%s", helpers.UUID(), handler.Filename)
				picturePath := filepath.Join(helpers.Storage, fileName)

				out, err := os.Create(picturePath)
				if err != nil {
					helpers.ErrorSimple(err)
					helpers.JSONResponse(w, errors.New("FAILED TO CREATE IMAGE FILE"), http.StatusInternalServerError)
					return
				}
				defer out.Close()

				_, err = out.Write(fileBytes)
				if err != nil {
					helpers.ErrorSimple(err)
					helpers.JSONResponse(w, errors.New("FAILED TO WRITE IMAGE CONTENT"), http.StatusInternalServerError)
					return
				}

				dataToUpdate["picture"] = fileName
			}

			champ.Update(dataToUpdate, "id", postID.Values["id"])
		}
		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusOK)
		global.Status = 200
		http.Redirect(w, r, "/", http.StatusNotFound)
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)
	posts := models.Post.Fetch()
	postID := posts[helpers.RandN(0, 1)]
	postToDelete := models.Post.Fetch([]interface{}{"id"}, []interface{}{postID.Values["id"]})

	if len(postToDelete) != 0 {

		for _, post := range postToDelete {
			post.Delete()
		}

		helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
	} else {
		http.Redirect(w, r, "/", http.StatusNotFound)
	}
}
