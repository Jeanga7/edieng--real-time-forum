package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/global"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, err error, context string, statusCode int) {
	global.W = w
	if err != nil {
		JSONResponse(w, err, statusCode)
	}
}

func JSONResponse(w http.ResponseWriter, err error, statusCode int, more ...interface{}) {
	global.W = w
	response := map[string]string{
		"value": err.Error(),
	}

	if len(more) > 0 {
		response["more"] = more[0].(string)
	}

	// Convertir la réponse en JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("%s: %v", err.Error(), err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Définir l'entête de la réponse
	w.Header().Set("Content-Type", "application/json")

	// Envoyer la réponse avec un code d'erreur
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func ErrorSimple(err error) {
	if err != nil {
		fmt.Println(err)
		var w http.ResponseWriter
		JSONResponse(w, err, http.StatusInternalServerError)
	}
}

func AllowedMethod(w http.ResponseWriter, allowed string, r *http.Request) {
	global.W = w
	if r.Method != allowed {
		ErrorHandler(w, errors.New("BAD REQUEST"), "AllowedMethod", http.StatusBadRequest)
		return
	}
}
