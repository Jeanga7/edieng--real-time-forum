package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func FetchDiscussions(w http.ResponseWriter, r *http.Request) {
	var discussions []models.Model
	params := r.URL.Query()

	if toShow := params.Get("discussions"); len(toShow) > 0 && toShow == "all" {
		discussions = models.Discussions.Fetch()
	}

	usersData := func() []map[string]interface{} {
		allValues := []map[string]interface{}{}
		for _, u := range discussions {
			var data = make(map[string]interface{})
			data["id"] = u.Values["id"]
			data["idSender"] = u.Values["idSender"]
			data["idReceiver"] = u.Values["idReceiver"]
			data["messages"] = u.Values["messages"]
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
