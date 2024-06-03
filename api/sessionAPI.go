package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func FetchSessions(w http.ResponseWriter, r *http.Request) {
	var Sessions []models.Model
	params := r.URL.Query()

	if toShow := params.Get("session"); len(toShow) > 0 && toShow == "all" {
		Sessions = models.Session.Fetch()
	}

	usersData := func() []map[string]interface{} {
		allValues := []map[string]interface{}{}
		for _, u := range Sessions {
			var data = make(map[string]interface{})
			data["id"] = u.Values["id"]
			data["idUser"] = u.Values["idUser"]
			data["expiryDate"] = u.Values["expiryDate"]
			data["creationDate"] = u.Values["creationDate"]
			allValues = append(allValues, data)
		}
		return allValues
	}()

	res, err := json.Marshal(usersData)
	helpers.ErrorSimple(err)

	// Header
	w.Header().Set("Content-Type", "application/json")

	// Response
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
