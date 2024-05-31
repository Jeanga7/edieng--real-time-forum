package controllers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var dest string

type Message struct {
	Mess  string `json:"mess"`
	Uname string `json:"uname"`
	Date  string `json:"date"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandlerWebs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	for {
		handleMessage(conn, r)
	}
}

func handleMessage(conn *websocket.Conn, r *http.Request) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}
	fmt.Println(string(msg))
	var message Message
	err = json.Unmarshal(msg, &message)
	if err != nil {
		dest = string(msg)
		sendPreviousMessages(conn, r, msg)
		return
	}

	IdDestinataire := GetIdByUsername(dest)
	_, _, _, IdEmetteur := FindUsersOnline(r)
	idUsers := []string{IdDestinataire, IdEmetteur.(string)}

	err = saveDiscussion(idUsers, message.Mess)
	if err != nil {
		log.Println("Error saving discussion:", err)
	}
}

func sendPreviousMessages(conn *websocket.Conn, r *http.Request, msg []byte) {
	IdDestinataire := GetIdByUsername(string(msg))
	_, _, _, IdEmetteur := FindUsersOnline(r)
	

	var myDiscussions []string
	discussion := models.Discussions.Fetch([]interface{}{"idSender"}, []interface{}{IdEmetteur})
	for _, v := range discussion {
		if v.Values["idReceiver"] == IdDestinataire {
			fmt.Println(v.Values["messages"])
			if v.Values["messages"].(string) != "" {
				var all = "[" + v.Values["messages"].(string) + v.Values["creationDate"].(string) + "]"
				fmt.Println(all)
				myDiscussions = append(myDiscussions, all)
			}
		}
	}
	if len(discussion) == 0 {
		return
	}

	
	baite := []byte(strings.Join(myDiscussions, ","))
	

	if len(myDiscussions) != 0 {
		if err := conn.WriteMessage(websocket.TextMessage, baite); err != nil {
			log.Println("Error sending previous messages:", err)
		}
	}
}


func saveDiscussion(idUsers []string, discussion string) error {
	
	fmt.Println("discussion", discussion)
	if discussion != "" {
		newDiscussion := models.Discussions

		newDiscussion.Values = map[string]interface{}{
			"idSender":   idUsers[1],
			"idReceiver": idUsers[0],
			"messages":   discussion,
		}
		newDiscussion.Create()
	}
	return nil
}

func GetIdByUsername(username string) string {
	conn := models.User.Fetch([]interface{}{"name"}, []interface{}{username})
	if len(conn) > 0 {
		return conn[0].Values["id"].(string)
	}
	return ""
}

func GetDiscussionId(User1, User2 string) interface{} {
	resIdDiscussion1, err := json.Marshal([]string{User1, User2})
	if err != nil {
		log.Println("Error marshalling discussion ID 1:", err)
		return ""
	}
	resIdDiscussion2, err := json.Marshal([]string{User2, User1})
	if err != nil {
		log.Println("Error marshalling discussion ID 2:", err)
		return ""
	}

	Diss1 := models.Discussions.Fetch([]interface{}{"idUsers"}, []interface{}{string(resIdDiscussion1)})
	Diss2 := models.Discussions.Fetch([]interface{}{"idUsers"}, []interface{}{string(resIdDiscussion2)})

	if len(Diss1) > 0 {
		return Diss1[0].Values["id"]
	}

	if len(Diss2) > 0 {
		return Diss2[0].Values["id"]
	}

	return ""
}
