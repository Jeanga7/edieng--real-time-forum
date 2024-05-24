package main

import (
	"forum/controllers"
	"net/http"
)

func Routes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./views/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	mux.HandleFunc("/", controllers.Index)

	// User routes
	mux.HandleFunc("/register", controllers.RegisterHandler)
	mux.HandleFunc("/login", controllers.LoginHandler)
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
}
