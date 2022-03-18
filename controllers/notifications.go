package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/notifications"

	"github.com/owncast/owncast/utils"

	log "github.com/sirupsen/logrus"
)

// RegisterForLiveNotifications will register a channel + destination to be
// notified when a stream goes live.
func RegisterForLiveNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != POST {
		WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	type request struct {
		// Channel is the notification channel (browser, sms, etc)
		Channel string `json:"channel"`
		// Destination is the target of the notification in the above channel.
		Destination string `json:"destination"`
	}

	decoder := json.NewDecoder(r.Body)
	var req request
	if err := decoder.Decode(&req); err != nil {
		log.Errorln(err)
		WriteSimpleResponse(w, false, "unable to register for notifications")
		return
	}

	// Make sure the requested channel is one we want to handle.
	validTypes := []string{notifications.BrowserPushNotification}
	_, validChannel := utils.FindInSlice(validTypes, req.Channel)
	if !validChannel {
		WriteSimpleResponse(w, false, "invalid notification channel: "+req.Channel)
		return
	}

	if err := notifications.AddNotification(req.Channel, req.Destination); err != nil {
		log.Errorln(err)
		WriteSimpleResponse(w, false, "unable to save notification")
		return
	}
}
