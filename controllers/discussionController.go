package controllers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"forum/models"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// func HandlerWebs(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Error upgrading connection:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		handleMessage(conn, r)
// 	}
// }

// func handleMessage(conn *websocket.Conn, r *http.Request) {
// 	_, msg, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Println("Error reading message:", err)
// 		return
// 	}
// 	fmt.Println(string(msg))
// 	var message Message
// 	err = json.Unmarshal(msg, &message)
// 	if err != nil {
// 		log.Println("Error unmarshalling message:", err)
// 		return
// 	}

// 	if message.Type == "fetch" {
// 		dest = message.Uname
// 		sendPreviousMessages(conn, r, message.Uname)
// 	} else if message.Type == "message" {
// 		IdDestinataire := GetIdByUsername(dest)
// 		nn, _, _, IdEmetteur := FindUsersOnline(r)
// 		fmt.Println("this is my name", nn)
// 		idUsers := []string{IdDestinataire, IdEmetteur.(string)}

// 		discussionId := GetDiscussionId(IdEmetteur.(string), IdDestinataire)
// 		fmt.Println("discussionId", discussionId)

// 		fmt.Println(message.Date)
// 		if discussionId == "" {
// 			err = saveDiscussion(idUsers, message.Mess, message.Uname, message.Date)
// 		} else {
// 			err = appendMessageToDiscussion(discussionId, message.Mess, message.Uname, message.Date)
// 		}

// 		if err != nil {
// 			log.Println("Error saving discussion:", err)
// 		}
// 	}
// }

// func sendPreviousMessages(conn *websocket.Conn, r *http.Request, username string) {
// 	IdDestinataire := GetIdByUsername(username)
// 	_, _, _, IdEmetteur := FindUsersOnline(r)

// 	var myDiscussions []Message

// 	// var Ids = GetDiscussionId(IdEmetteur.(string),IdDestinataire)
// 	// for _, v := range Ids.(string) {

// 	// }
// 	discussion := models.Discussions.Fetch([]interface{}{"idSender"}, []interface{}{IdEmetteur})
// 	for _, v := range discussion {
// 		if v.Values["idReceiver"] == IdDestinataire {
// 			if v.Values["messages"].(string) != "" {
// 				var messages []Message
// 				err := json.Unmarshal([]byte(v.Values["messages"].(string)), &messages)
// 				if err != nil {
// 					log.Println("Error unmarshalling messages:", err)
// 					continue
// 				}

// 				myDiscussions = append(myDiscussions, messages...)
// 			}
// 		}
// 	}

// 	if len(myDiscussions) > 0 {
// 		historyMessage := struct {
// 			Type     string    `json:"type"`
// 			Messages []Message `json:"messages"`
// 		}{
// 			Type:     "history",
// 			Messages: myDiscussions,
// 		}

// 		err := conn.WriteJSON(historyMessage)
// 		if err != nil {
// 			log.Println("Error sending previous messages:", err)
// 		}
// 	}
// }

// func saveDiscussion(idUsers []string, discussion string, username string, date string) error {
// 	if discussion != "" {
// 		var messages []Message

// 		newMessage := Message{
// 			Mess:  discussion,
// 			Uname: username,
// 			Date:  date,
// 		}
// 		messages = append(messages, newMessage)

// 		messagesJSON, err := json.Marshal(messages)
// 		if err != nil {
// 			return err
// 		}

// 		fmt.Println("test", string(messagesJSON))

// 		newDiscussion := models.Discussions
// 		newDiscussion.Values = map[string]interface{}{
// 			"idSender":   idUsers[1],
// 			"idReceiver": idUsers[0],
// 			"messages":   string(messagesJSON),
// 		}
// 		newDiscussion.Create()
// 	}
// 	return nil
// }

// func appendMessageToDiscussion(discussionId interface{}, message string, username string, date string) error {
// 	discussion := models.Discussions.Fetch([]interface{}{"id"}, []interface{}{discussionId})
// 	if len(discussion) > 0 {
// 		var messages []Message
// 		err := json.Unmarshal([]byte(discussion[0].Values["messages"].(string)), &messages)
// 		if err != nil {
// 			return err
// 		}

// 		newMessage := Message{
// 			Mess:  message,
// 			Uname: username,
// 			Date:  date,
// 		}
// 		messages = append(messages, newMessage)

// 		messagesJSON, err := json.Marshal(messages)
// 		if err != nil {
// 			return err
// 		}

// 		data := map[string]interface{}{"messages": string(messagesJSON)}
// 		models.Discussions.Update(data, "id", discussionId)
// 	}
// 	return nil
// }

// func GetIdByUsername(username string) string {
// 	conn := models.User.Fetch([]interface{}{"name"}, []interface{}{username})
// 	if len(conn) > 0 {
// 		return conn[0].Values["id"].(string)
// 	}
// 	return ""
// }

// func GetDiscussionId(User1, User2 string) interface{} {
// 	discussions := models.Discussions.Fetch()
// 	for _, discussion := range discussions {
// 		idSender := discussion.Values["idSender"].(string)
// 		idReceiver := discussion.Values["idReceiver"].(string)

// 		if (idSender == User1 && idReceiver == User2) || (idSender == User2 && idReceiver == User1) {
// 			return discussion.Values["id"]
// 		}
// 	}

// 	return ""
// }

// /////////////////////////////////////////////////////////////
