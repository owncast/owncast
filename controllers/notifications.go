package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/notifications"
	"github.com/owncast/owncast/notifications/mailjet"

	"github.com/owncast/owncast/utils"

	log "github.com/sirupsen/logrus"
)

// RegisterForEmailNotifications will register a single email address with
// an email list, creating the list if necessary.
func RegisterForEmailNotifications(w http.ResponseWriter, r *http.Request) {
	type request struct {
		EmailAddress string `json:"emailAddress"`
	}

	emailConfig := data.GetMailjetConfiguration()
	if !emailConfig.Enabled {
		WriteSimpleResponse(w, false, "email notifications are not enabled")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req request
	if err := decoder.Decode(&req); err != nil {
		log.Errorln(err)
		WriteSimpleResponse(w, false, "unable to register for email notifications")
		return
	}

	m := mailjet.New(emailConfig.Username, emailConfig.Password)

	// If we have not previously created an email list for Owncast then create
	// a new one now, and add the requested email address to it.
	if emailConfig.ListID == 0 {
		listAddress, listID, err := m.CreateListAndAddAddress(req.EmailAddress)
		if err != nil {
			log.Errorln(err)
			WriteSimpleResponse(w, false, "unable to register address for notifications")
			return
		}
		emailConfig.ListAddress = listAddress
		emailConfig.ListID = listID
		if err := data.SetMailjetConfiguration(emailConfig); err != nil {
			log.Errorln(err)
			WriteSimpleResponse(w, false, "error in saving email configuration")
			return
		}
	} else {
		if err := m.AddEmailToList(req.EmailAddress, emailConfig.ListID); err != nil {
			log.Errorln(err)
			WriteSimpleResponse(w, false, "error in adding email address to list")
			return
		}
	}

	WriteSimpleResponse(w, true, "added")
}

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
	validTypes := []string{notifications.BrowserPushNotification, notifications.TextMessageNotification}
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
