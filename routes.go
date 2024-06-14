package main

import (
	"forum/api"
	"forum/controllers"
	"net/http"
)

func Routes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./views/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	mux.HandleFunc("/messages", controllers.HandleMessages)
	go controllers.HandleTheMessages()
	mux.HandleFunc("/", controllers.Index)

	// User routes
	mux.HandleFunc("/register", controllers.RegisterHandler)
	mux.HandleFunc("/login", controllers.LoginHandler)
	mux.HandleFunc("/start", controllers.NewLogin)
	mux.HandleFunc("/logout", controllers.LogoutHandler)

	// Posts routes
	mux.HandleFunc("/create-post", controllers.CreatePost)
	mux.HandleFunc("/edit-post", controllers.EditPost)
	mux.HandleFunc("/delete-post", controllers.DeletePost)

	// Comments routes
	mux.HandleFunc("/create-comment", controllers.CreateComment)

	// Category routes
	// mux.HandleFunc("/add-category", controllers.AddCategory)

	// Like routes
	mux.HandleFunc("/like", controllers.Like)

	// Dislike routes
	mux.HandleFunc("/dislike", controllers.Dislike)

	// API Routes
	mux.HandleFunc("/api/users", api.FetchUsers)

	mux.HandleFunc("/api/auth", api.ValidateFields)

	mux.HandleFunc("/api/posts", api.FetchPosts)

	mux.HandleFunc("/api/category", api.FetchCategory)

	mux.HandleFunc("/api/categories-list", api.CategoriesList)

	mux.HandleFunc("/api/session", api.FetchSessions)

	mux.HandleFunc("/api/comments", api.FetchComments)

	mux.HandleFunc("/api/discussions", api.FetchDiscussions)

	mux.HandleFunc("/api/likes", api.FetchLike)

	mux.HandleFunc("/api/dislikes", api.FetchDislike)

	mux.HandleFunc("/api/all-users", api.GetAllSortedUsers)
}
