package controllers

import (
	"encoding/json"
	"forum/models"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

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

	var message Message
	err = json.Unmarshal(msg, &message)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		sendPreviousMessages(conn, r, msg)
		return
	}

	IdDestinataire := GetIdByUsername(message.Uname)
	_, _, _, IdEmetteur := FindUsersOnline(r)
	idUsers := []string{IdDestinataire, IdEmetteur.(string)}

	discussion, err := getDiscussionMessages(IdDestinataire, IdEmetteur.(string))
	if err != nil {
		log.Println("Error getting discussion messages:", err)
		return
	}

	discussion = append(discussion, string(msg))

	err = saveDiscussion(idUsers, discussion)
	if err != nil {
		log.Println("Error saving discussion:", err)
	}
}

func sendPreviousMessages(conn *websocket.Conn, r *http.Request, msg []byte) {
	IdDestinataire := GetIdByUsername(string(msg))
	_, _, _, IdEmetteur := FindUsersOnline(r)
	idUsers := []string{IdDestinataire, IdEmetteur.(string)}
	IdUserKey, err := json.Marshal(idUsers)
	if err != nil {
		log.Println("Error marshalling idUsers:", err)
		return
	}

	discussion := models.Discussions.Fetch([]interface{}{"idUsers"}, []interface{}{string(IdUserKey)})
	if len(discussion) == 0 {
		return
	}

	result := discussion[0].Values["messages"]
	toSend := strings.ReplaceAll(result.(string), "\\", "")
	toWebs, err := json.Marshal(toSend)
	if err != nil {
		log.Println("Error marshalling previous messages:", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, toWebs); err != nil {
		log.Println("Error sending previous messages:", err)
	}
}

func getDiscussionMessages(User1, User2 string) ([]string, error) {
	IdDisc := GetDiscussionId(User1, User2)
	if IdDisc == "" {
		return []string{}, nil
	}

	existingDiscussion := models.Discussions.Fetch([]interface{}{"id"}, []interface{}{IdDisc.(string)})
	if len(existingDiscussion) == 0 {
		return []string{}, nil
	}

	res := existingDiscussion[0].Values["messages"]
	existingDiscussion[0].Delete()

	var discussion []string
	err := json.Unmarshal([]byte(res.(string)), &discussion)
	if err != nil {
		return nil, err
	}

	return discussion, nil
}

func saveDiscussion(idUsers []string, discussion []string) error {
	res, err := json.Marshal(discussion)
	if err != nil {
		return err
	}

	resIdUsers, err := json.Marshal(idUsers)
	if err != nil {
		return err
	}
	newDiscussion := models.Discussions

	newDiscussion.Values = map[string]interface{}{
		"idUsers":  string(resIdUsers),
		"messages": string(res),
	}
	newDiscussion.Create()
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
