package main

import (
	"fmt"
	"forum/api"
	"forum/controllers"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Routes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./views/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	mux.HandleFunc("/messages", HandleMessages)
	go HandletheMessages()
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
}

var clients = make(map[*websocket.Conn]string)

var message = make(chan Message)

var dest string

type Message struct {
	Mess  string `json:"mess"`
	Uname string `json:"uname"`
	Date  string `json:"date"`
	Type  string `json:"type"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleMessages(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	var initialMessage Message
	err = conn.ReadJSON(&initialMessage)
	if err != nil {
		log.Println("Error reading initial message:", err)
		return
	}

	name, _, _, _ := controllers.FindUsersOnline(r)
	fmt.Println(name)
	clients[conn] = name

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, conn)
			break
		}
		message <- msg
	}
}

func HandletheMessages() {
	for {
		msg := <-message
		for client, username := range clients {
			if username == msg.Uname {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
