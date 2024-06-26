package controllers

import (
	"encoding/json"
	"fmt"
	"forum/helpers"
	"forum/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var connectedNames = make(map[string]*websocket.Conn)
var message = make(chan Message)

type Message struct {
	MessageContent string `json:"message"`
	Sender         string `json:"sender"`
	Destinataire   string `json:"destinataire"`
	Date           string `json:"date"`
	Notif          Notification
	Type           string `json:"type"`
}

type ClientMessage struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Signature string `json:"signature"`
	Time      string `json:"time"`
}

type Notification struct {
	Sender              string `json:"sender"`
	Destinataire        string `json:"destinataire"`
	NombreNotifications string `json:"notification"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleMessages(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	ReadLoop(r, conn)
}

func ReadLoop(r *http.Request, conn *websocket.Conn) {
	for {
		// Read message
		var msg ClientMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			name, _, _, _ := FindUsersOnline(r)
			delete(connectedNames, name)
			go Broadcasting(name, "offline")
			break
		}

		// Find message author
		resFrom := models.Session.Fetch([]interface{}{"id"}, []interface{}{msg.Signature})
		resTo := models.User.Fetch([]interface{}{"nickname"}, []interface{}{msg.To})

		if len(resFrom) > 0 {
			userFound := models.User.Fetch([]interface{}{"id"}, []interface{}{resFrom[0].Values["idUser"]})[0]
			senderId, senderNickname := userFound.Values["id"].(string), userFound.Values["nickname"].(string)

			// Check if connected user
			if connectedNames[senderNickname] != nil && len(resTo) > 0 {
				dest := resTo[0]
				destId, destNickname := dest.Values["id"].(string), dest.Values["nickname"].(string)

				go SaveAndSend(conn, senderId, senderNickname, destId, destNickname, msg)

				// Print message on server side
				// fmt.Printf("Received message from %s: '%s'\n", senderNickname, string(msg.Message))
			} else if string(msg.To) == "@system" {
				connectedNames[senderNickname] = conn // Register the senderNickname
				fmt.Println(senderNickname + " connected")

				go Broadcasting(senderNickname, "online")
			}
			// {"to":"@system","type":"join","message":"","signature":"035fe39c-b3db-4e01-8062-2a4af9a1f66f"}
		}
	}
}

func SaveAndSend(conn *websocket.Conn, senderId, senderNickname, destId, destNickname string, msg ClientMessage) {
	if msg.Type == "message" {
		// Save
		err := StoreMessage(senderId, senderNickname, destId, destNickname, msg.Message)
		if err != nil {
			res := ClientMessage{
				From:      "@server",
				To:        senderNickname,
				Message:   err.Error(),
				Type:      "error",
				Signature: "@server",
			}
			txt, err := json.Marshal(res)
			helpers.ErrorSimple(err)
			conn.WriteJSON(string(txt))
			return
		}

		// Save notification
		notifs := models.Notification.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{senderId, destId})
		if len(notifs) == 0 {
			// Create notification
			CreateNotification(senderId, destId)
		} else {
			// Update notification
			notifId := notifs[0].Values["id"].(string)
			notifNb := notifs[0].Values["notifications"].(int64)
			UpdateNotif(notifId, senderId, destId, notifNb+1)
		}
	}

	// Send
	if connectedNames[destNickname] != nil {
		res := ClientMessage{
			From:      senderNickname,
			To:        destNickname,
			Message:   msg.Message,
			Type:      msg.Type,
			Signature: msg.Signature,
			Time:      helpers.Now(),
		}
		txt, err := json.Marshal(res)
		helpers.ErrorSimple(err)

		err = connectedNames[destNickname].WriteJSON(string(txt))
		helpers.ErrorSimple(err)
	}
}

func Broadcasting(user, typeMsg string) {
	// Create message
	res := ClientMessage{
		From:      "@server",
		To:        "",
		Message:   user,
		Type:      typeMsg,
		Signature: "@server",
		Time:      helpers.Now(),
	}

	// Broadcast
	for n := range connectedNames {
		if n != user {
			res.To = n
			txt, err := json.Marshal(res)
			helpers.ErrorSimple(err)
			err = connectedNames[n].WriteJSON(string(txt))
			helpers.ErrorSimple(err)
		}
	}
}

func StoreMessage(senderId, senderNickname, destId, destNickname, message string) error {
	discussions := models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{senderId, destId})
	if len(discussions) == 0 {
		discussions = models.Discussions.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{destId, senderId})
	}

	if len(discussions) == 0 {
		return SaveDiscussion(message, senderId, destId, helpers.Now())
	} else {
		discussionId := GetDiscussionId(senderId, destId)
		return AppendMessageToDiscussion(discussionId, message, senderNickname, destNickname, helpers.Now())
	}
}

// func HandleMessages(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Error upgrading connection:", err)
// 		return
// 	}
// 	defer conn.Close()
// 	var initialMessage Message
// 	err = conn.ReadJSON(&initialMessage)

// 	name, _, _, _ := FindUsersOnline(r)

// 	if err != nil {
// 		log.Println("Error reading initial message:", err)
// 		return
// 	}

// 	clients[conn] = name
// 	for {

// 		var msg Message
// 		err := conn.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("Error: %v", err)
// 			delete(clients, conn)
// 			break
// 		}

// 		discussionId := GetDiscussionId(msg.Sender, msg.Destinataire)

// 		if discussionId == "" {
// 			err = SaveDiscussion(msg.MessageContent, msg.Sender, msg.Destinataire, msg.Date)
// 		} else {
// 			err = AppendMessageToDiscussion(discussionId, msg.MessageContent, msg.Sender, msg.Destinataire, msg.Date)
// 		}

// 		notificationId := GetNotificationId(msg.Sender, msg.Destinataire)
// 		if notificationId == "" {
// 			CreateNotification(msg.Sender, msg.Destinataire)
// 		} else if msg.Type == "fetch" {
// 			notificationId := GetNotificationId(msg.Destinataire, name)
// 			notification := models.Notification.Fetch([]interface{}{"id"}, []interface{}{notificationId})
// 			SendPreviousMessages(conn, msg.Sender, msg.Destinataire)
// 			if len(notification) > 0 {
// 				UpdateNotif(notificationId, msg.Sender, name, 0)
// 			}
// 		} else {
// 			if msg.Type != "typing" {
// 				if msg.MessageContent != "" {

// 					notification := models.Notification.Fetch([]interface{}{"id"}, []interface{}{notificationId})
// 					if len(notification) > 0 {
// 						currentCount, err := strconv.Atoi(notification[0].Values["notifications"].(string))
// 						if err != nil {
// 							return
// 						}

// 						number := currentCount + 1
// 						UpdateNotif(notificationId, msg.Sender, name, number)
// 					}
// 				}
// 			}
// 		}

// 		if err != nil {
// 			log.Println("Error saving discussion:", err)
// 		}
// 		// SendPreviousMessages(conn, msg.Sender, msg.Destinataire)
// 		message <- msg
// 	}
// }

func HandleTheMessages() {
	for {
		msg := <-message
		// fmt.Println("hello", msg)

		for client, nickname := range clients {
			if nickname == msg.Destinataire {
				if msg.Type != "typing" {
					if msg.MessageContent != "" {
						notificationId := GetNotificationId(msg.Sender, msg.Destinataire)
						notification := models.Notification.Fetch([]interface{}{"id"}, []interface{}{notificationId})
						if len(notification) > 0 {
							currentCount, err := strconv.Atoi(notification[0].Values["notifications"].(string))
							if err != nil {
								return
							}
							msg.Notif.NombreNotifications = strconv.Itoa(currentCount)
						}
					}
				}
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

func SendPreviousMessages(conn *websocket.Conn, sender, receiver string) {
	var myDiscussions []Message

	discussion := models.Discussions.Fetch()
	for _, v := range discussion {
		if v.Values["idSender"] == sender && v.Values["idReceiver"] == receiver ||
			v.Values["idSender"] == receiver && v.Values["idReceiver"] == sender {
			if v.Values["messages"].(string) != "" {
				var messages []Message
				err := json.Unmarshal([]byte(v.Values["messages"].(string)), &messages)
				if err != nil {
					log.Println("Error unmarshalling messages:", err)
					continue
				}
				myDiscussions = append(myDiscussions, messages...)
			}
		}
	}

	if len(myDiscussions) > 0 {
		historyMessage := struct {
			Type     string    `json:"type"`
			Messages []Message `json:"messages"`
		}{
			Type:     "history",
			Messages: myDiscussions,
		}
		err := conn.WriteJSON(historyMessage)
		if err != nil {
			log.Println("Error sending previous messages:", err)
		}
	}
}

func SaveDiscussion(content, sender, receiver, date string) error {
	if content != "" {
		var messages []Message

		usersFound1 := models.User.Fetch([]interface{}{"id"}, []interface{}{sender})
		usersFound2 := models.User.Fetch([]interface{}{"id"}, []interface{}{receiver})

		var chatSender, chatReceiver string

		if len(usersFound1) == 1 && len(usersFound2) == 1 {
			chatSender = (usersFound1[0].Values["nickname"]).(string)
			chatReceiver = (usersFound2[0].Values["nickname"]).(string)
		}

		newMessage := Message{
			MessageContent: content,
			Sender:         chatSender,
			Destinataire:   chatReceiver,
			Date:           date,
		}
		messages = append(messages, newMessage)

		messagesJSON, err := json.Marshal(messages)
		if err != nil {
			return err
		}

		newDiscussion := models.Discussions
		newDiscussion.Values = map[string]interface{}{
			"idSender":   sender,
			"idReceiver": receiver,
			"messages":   string(messagesJSON),
		}
		newDiscussion.Create()
	}
	return nil
}

func CreateNotification(sender string, receiver string) {
	notif := models.Notification
	notif.Values = map[string]interface{}{
		"idSender":      sender,
		"idReceiver":    receiver,
		"notifications": 1,
	}
	notif.Create()
}

func UpdateNotif(notifId interface{}, sender, receiver string, number int64) error {
	data := map[string]interface{}{"notifications": number}
	models.Notification.Update(data, "id", notifId)
	return nil
}

func AppendMessageToDiscussion(discussionId interface{}, message, sender, receiver, date string) error {
	if message != "" {
		discussion := models.Discussions.Fetch([]interface{}{"id"}, []interface{}{discussionId})
		if len(discussion) > 0 {
			var messages []Message
			err := json.Unmarshal([]byte(discussion[0].Values["messages"].(string)), &messages)
			if err != nil {
				return err
			}

			newMessage := Message{
				MessageContent: message,
				Sender:         sender,
				Destinataire:   receiver,
				Date:           date,
			}
			messages = append(messages, newMessage)

			messagesJSON, err := json.Marshal(messages)
			if err != nil {
				return err
			}

			data := map[string]interface{}{"messages": string(messagesJSON)}
			models.Discussions.Update(data, "id", discussionId)
		}
	}
	return nil
}
func GetDiscussionId(User1, User2 string) interface{} {
	discussions := models.Discussions.Fetch()
	for _, discussion := range discussions {
		idSender := discussion.Values["idSender"].(string)
		idReceiver := discussion.Values["idReceiver"].(string)

		if (idSender == User1 && idReceiver == User2) || (idSender == User2 && idReceiver == User1) {
			return discussion.Values["id"]
		}
	}
	return ""
}

func GetNotificationId(User1, User2 string) interface{} {
	Notifications := models.Notification.Fetch()
	for _, notif := range Notifications {
		idSender := notif.Values["idSender"].(string)
		idReceiver := notif.Values["idReceiver"].(string)

		if idSender == User1 && idReceiver == User2 {
			return notif.Values["id"]
		}
	}
	return ""
}

func BroadcastNewUser(nickname string) {
	fmt.Println("nickname", nickname)
	newUserMessage := Message{
		Type:   "new_user",
		Sender: nickname,
	}
	for client := range clients {
		err := client.WriteJSON(newUserMessage)
		if err != nil {
			log.Printf("Error broadcasting new user to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
