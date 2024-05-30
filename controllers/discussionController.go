package controllers

import (
	"encoding/json"
	"fmt"
	"forum/helpers"
	"forum/models"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// type discussion struct {
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

		// fmt.Println("message:", string(msg))

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			// log.Println("Erreur JSON:", err)
			IdDestinataire := GetIdByUsername(string(msg))
			_, _, _, IdEmetteur := FindUsersOnline(r)
			idUsers := []string{IdDestinataire, IdEmetteur.(string)}
			IdUserKey, err := json.Marshal(idUsers)
			helpers.ErrorSimple(err)

			// fmt.Println("IdUserKey:", string(IdUserKey))
			var Discontent = models.Discussions.Fetch([]interface{}{"idUsers"}, []interface{}{string(IdUserKey)})
			var result = Discontent[0].Values["messages"]

			// var temp Message
			// err = json.Unmarshal([]byte(result.(string)), &temp)
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }
			// fmt.Println(temp)
			var to_send = strings.ReplaceAll(result.(string), "\\", "")
			to_webs, err := json.Marshal(to_send)
			if err != nil {
				log.Println("Error encoding scoreboard message:", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, to_webs); err != nil {
				log.Println("Error sending scoreboard message:", err)
				continue
			}
			continue
		}

		IdDestinataire := GetIdByUsername(message.Uname)
		_, _, _, IdEmetteur := FindUsersOnline(r)
		idUsers := []string{IdDestinataire, IdEmetteur.(string)}

		var discussion []string
		IdDisc := GetDiscussionId(IdDestinataire, IdEmetteur.(string))
		fmt.Println("IdDisc", IdDisc.(string))

		if IdDisc != "" {
			existingDiscussion := models.Discussions.Fetch([]interface{}{"id"}, []interface{}{IdDisc.(string)})
			if len(existingDiscussion) > 0 {
				res := existingDiscussion[0].Values["messages"]
				existingDiscussion[0].Delete()
				err = json.Unmarshal([]byte(res.(string)), &discussion)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}

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

func GetDiscussionId(User1, User2 string) interface{} {
	resIdDiscussion1, err := json.Marshal([]string{User1, User2})
	helpers.ErrorSimple(err)
	resIdDiscussion2, err := json.Marshal([]string{User2, User1})
	helpers.ErrorSimple(err)

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
