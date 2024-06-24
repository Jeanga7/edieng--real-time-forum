package api

import (
	"errors"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func ResetNotifications(w http.ResponseWriter, r *http.Request) {
	helpers.AllowedMethod(w, "POST", r)

	sender, dest := r.FormValue("sender"), r.FormValue("dest")

	if len(sender) > 0 && len(dest) > 0 {
		senders := models.User.Fetch([]interface{}{"nickname"}, []interface{}{sender})
		dests := models.User.Fetch([]interface{}{"nickname"}, []interface{}{dest})
		if len(senders) > 0 && len(dests) > 0 {
			idSender, idReceiver := senders[0].Values["id"].(string), dests[0].Values["id"].(string)

			notifs := models.Notification.Fetch([]interface{}{"idSender", "idReceiver"}, []interface{}{idSender, idReceiver})
			if len(notifs) > 0 {
				data := map[string]interface{}{"notifications": 0}
				models.Notification.Update(data, "id", notifs[0].Values["id"])
			}
			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusOK)
			return
		} else {
			helpers.JSONResponse(w, errors.New("BAD INFORMATIONS"), http.StatusBadRequest)
		}
	} else {
		helpers.JSONResponse(w, errors.New("MISSING INFORMATIONS"), http.StatusBadRequest)
	}
}
