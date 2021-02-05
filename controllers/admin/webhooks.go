package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

type deleteWebhookRequest struct {
	ID int `json:"id"`
}

type createWebhookRequest struct {
	URL    string             `json:"url"`
	Events []models.EventType `json:"events"`
}

// CreateWebhook will add a single webhook.
func CreateWebhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request createWebhookRequest
	if err := decoder.Decode(&request); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	// Verify all the scopes provided are valid
	if !models.HasValidEvents(request.Events) {
		controllers.BadRequestHandler(w, errors.New("one or more invalid event provided"))
		return
	}

	newWebhookID, err := data.InsertWebhook(request.URL, request.Events)
	if err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteResponse(w, models.Webhook{
		ID:        newWebhookID,
		Url:       request.URL,
		Events:    request.Events,
		Timestamp: time.Now(),
		LastUsed:  nil,
	})
}

// GetWebhooks will return all webhooks.
func GetWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := data.GetWebhooks()
	if err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteResponse(w, webhooks)
}

// DeleteWebhook will delete a single webhook.
func DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request deleteWebhookRequest
	if err := decoder.Decode(&request); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	if err := data.DeleteWebhook(request.ID); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteSimpleResponse(w, true, "deleted webhook")
}
