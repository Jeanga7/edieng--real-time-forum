package controllers

import (
	"encoding/json"
	"fmt"
	"forum/helpers"
	"forum/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// type Message struct {
// 	Id       int
// 	IdUsers  string
// 	Messages string
// }

type Message struct {
	Mess  string `json:"mess"`
	Uname string `json:"uname"`
	Date  string `json:"date"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandlerWebs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("message:", string(msg))

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Erreur JSON:", err)
			continue
		}

		IdDestinataire := GetIdByUsername(message.Uname)
		_, _, _, IdEmetteur := FindUsersOnline(r)
		idUsers := []string{IdDestinataire, IdEmetteur.(string)}
		// fmt.Println("Mess:", message.Mess)
		// fmt.Println("IdDestinataire:", IdDestinataire)
		// fmt.Println("IdEmetteur:", IdEmetteur)

		var discussion []string
		// existingDiscussion, err := GetDiscussion("toi", "moi")
		if err != nil {
			log.Println(err)
			return
		}
		// if existingDiscussion != nil {
		// err = json.Unmarshal([]byte(existingDiscussion.Messages), &discussion)
		// 	if err != nil {
		// 		log.Println(err)
		// 		return
		// 	}
		// }

		discussion = append(discussion, string(msg))

		res, err := json.Marshal(discussion)
		helpers.ErrorSimple(err)

		resIdUsers, err := json.Marshal(idUsers)
		helpers.ErrorSimple(err)
		// // On vérifie si la discussion existe dans la base de données
		// // Si oui on met à jour la valeur de allDiscussion avec les messages déjà existant
		newDiscussion := models.Discussions

		newDiscussion.Values = map[string]interface{}{
			"idUsers":  string(resIdUsers),
			"messages": string(res),
		}
		newDiscussion.Create()
	}
}

func GetIdByUsername(username string) (id string) {
	conn := models.User.Fetch([]interface{}{"name"}, []interface{}{username})
	if len(conn) > 0 {
		id = conn[0].Values["id"].(string)
	}
	return id
}

func GetDiscussion(User1, User2 string) {

}
