package controllers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var SenderName, DestName string

type Message struct {
	MessageContent string `json:"message"`
	Sender         string `json:"sender"`
	Destinataire   string `json:"destinataire"`
	Date           string `json:"date"`
	Type           string `json:"type"`
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
		handleMessage(conn)
	}
}

func handleMessage(conn *websocket.Conn) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}
	// fmt.Println("msg", string(msg))
	var message Message
	err = json.Unmarshal(msg, &message)
	if err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}

	fmt.Println("message", message)

	if message.Type == "fetch" {
		SenderName = message.Sender
		DestName = message.Destinataire
		SendPreviousMessages(conn, message.Sender, message.Destinataire)
	} else if message.Type == "message" {

		discussionId := GetDiscussionId(SenderName, DestName)

		if discussionId == "" {
			err = SaveDiscussion(message.MessageContent, message.Sender, message.Destinataire, message.Date)
		} else {
			err = AppendMessageToDiscussion(discussionId, message.MessageContent, message.Sender, message.Destinataire, message.Date)
		}

		if err != nil {
			log.Println("Error saving discussion:", err)
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
