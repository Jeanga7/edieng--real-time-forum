package controllers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var message = make(chan Message)

type Message struct {
	MessageContent string `json:"message"`
	Sender         string `json:"sender"`
	Destinataire   string `json:"destinataire"`
	Date           string `json:"date"`
	Notif          Notification
	Type           string `json:"type"`
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
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
	var initialMessage Message
	err = conn.ReadJSON(&initialMessage)

	name, _, _, _ := FindUsersOnline(r)

	if err != nil {
		log.Println("Error reading initial message:", err)
		return
	}
	if initialMessage.Type == "notification" {
		fmt.Println("dès le début......")
	}

	clients[conn] = name
	for {

		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, conn)
			break
		}
		discussionId := GetDiscussionId(msg.Sender, msg.Destinataire)

		if discussionId == "" {
			err = SaveDiscussion(msg.MessageContent, msg.Sender, msg.Destinataire, msg.Date)
		} else {
			err = AppendMessageToDiscussion(discussionId, msg.MessageContent, msg.Sender, msg.Destinataire, msg.Date)
		}

		notificationId := GetNotificationId(msg.Sender, msg.Destinataire)
		fmt.Println(msg)
		if notificationId == "" {
			CreateNotification(msg.Sender, msg.Destinataire)
		} else if msg.Type == "fetch" {
			fmt.Println(msg.Destinataire, name)
			notificationId := GetNotificationId(msg.Destinataire, name)
			fmt.Println(notificationId)
			notification := models.Notification.Fetch([]interface{}{"id"}, []interface{}{notificationId})
			fmt.Println("ù*******************", notification)
			SendPreviousMessages(conn, msg.Sender, msg.Destinataire)
			if len(notification) > 0 {
				AddNotif(notificationId, msg.Sender, name, 0)
			}
		} else {
			notification := models.Notification.Fetch([]interface{}{"id"}, []interface{}{notificationId})
			if len(notification) > 0 {
				currentCount, err := strconv.Atoi(notification[0].Values["notifications"].(string))
				if err != nil {
					return
				}

				number := currentCount + 1
				AddNotif(notificationId, msg.Sender, name, number)
			}
		}

		if err != nil {
			log.Println("Error saving discussion:", err)
		}
		// SendPreviousMessages(conn, msg.Sender, msg.Destinataire)
		message <- msg
	}
}

func HandleTheMessages() {
	for {
		msg := <-message
		for client, username := range clients {
			if username == msg.Destinataire {
				notificationId := GetNotificationId(msg.Sender, msg.Destinataire)
				notification := models.Notification.Fetch([]interface{}{"id"}, []interface{}{notificationId})
				if len(notification) > 0 {
					currentCount, err := strconv.Atoi(notification[0].Values["notifications"].(string))
					if err != nil {
						return
					}
					msg.Notif.NombreNotifications = strconv.Itoa(currentCount)
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

		newMessage := Message{
			MessageContent: content,
			Sender:         sender,
			Destinataire:   receiver,
			Date:           date,
		}
		messages = append(messages, newMessage)

		messagesJSON, err := json.Marshal(messages)
		if err != nil {
			return err
		}

		// fmt.Println("test", string(messagesJSON))

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
		"notifications": "1",
	}
	notif.Create()
}

func AddNotif(notifId interface{}, sender, receiver string, number int) error {

	data := map[string]interface{}{"notifications": strconv.Itoa(number)}
	models.Notification.Update(data, "id", notifId)

	return nil
}

func AppendMessageToDiscussion(discussionId interface{}, message, sender, receiver, date string) error {
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
