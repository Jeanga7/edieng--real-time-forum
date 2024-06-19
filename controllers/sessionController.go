package controllers

import (
	"forum/global"
	"forum/models"
	"net/http"
	"time"
)

func CreateSession(w http.ResponseWriter, r *http.Request, idUser string) models.Model {
	global.W = w
	IdUser = idUser
	// Define the expiration date of the session (e.g., 1 hour)
	expiryDate := time.Now().Add(time.Hour)

	oldSessions := models.Session.Fetch([]interface{}{"idUser"}, []interface{}{idUser})
	if len(oldSessions) > 0 {
		oldSessions[0].Delete()
	}

	// Create a new session in the database
	session := models.Session
	session.Values = map[string]interface{}{
		"idUser":     idUser, // Ensure the email is used as the user identifier
		"expiryDate": expiryDate,
	}
	newSession := session.Create()
	IdSession = newSession.Values["id"].(string)

	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    IdSession,
		Path:     "/",
		Expires:  expiryDate,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	return newSession
}

func DeleteSession(w http.ResponseWriter, r *http.Request, idSession string) {
	global.W = w
	models.Session.Delete("id", "=", idSession)
	IdUser = ""

	sessionCookie := &http.Cookie{
		Name:     "sessionID",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, sessionCookie)
}
