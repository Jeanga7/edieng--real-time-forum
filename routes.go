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
	mux.HandleFunc("/", controllers.NewLogin)

	// User routes
	// mux.HandleFunc("/start", controllers.NewLogin)
	mux.HandleFunc("/logout", controllers.LogoutHandler)

	// Like routes
	mux.HandleFunc("/reaction", controllers.React)

	// API Routes
	mux.HandleFunc("/api/users", api.FetchUsers)

	mux.HandleFunc("/api/auth", api.ValidateFields)

	mux.HandleFunc("/api/posts", api.FetchPosts)

	mux.HandleFunc("/api/category", api.FetchCategory)

	mux.HandleFunc("/api/categories-list", api.CategoriesList)

	mux.HandleFunc("/api/session", api.FetchSessions)

	mux.HandleFunc("/api/comments", api.FetchComments)

	mux.HandleFunc("/api/create-comment", api.CreateComment)

	mux.HandleFunc("/api/discussions", api.HandleDiscussions)

	mux.HandleFunc("/api/likes", api.FetchLike)

	mux.HandleFunc("/api/dislikes", api.FetchDislike)

	mux.HandleFunc("/api/all-users", api.GetAllSortedUsers)

	mux.HandleFunc("/api/reset-notifications", api.ResetNotifications)
}
