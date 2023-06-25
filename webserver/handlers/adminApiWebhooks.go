package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage"
	"github.com/owncast/owncast/webserver/responses"
)

type deleteWebhookRequest struct {
	ID int `json:"id"`
}

type createWebhookRequest struct {
	URL    string             `json:"url"`
	Events []models.EventType `json:"events"`
}

// CreateWebhook will add a single webhook.
func (h *Handlers) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	webhookRepository := storage.GetWebhookRepository()

	decoder := json.NewDecoder(r.Body)
	var request createWebhookRequest
	if err := decoder.Decode(&request); err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	// Verify all the scopes provided are valid
	if !models.HasValidEvents(request.Events) {
		responses.BadRequestHandler(w, errors.New("one or more invalid event provided"))
		return
	}

	newWebhookID, err := webhookRepository.InsertWebhook(request.URL, request.Events)
	if err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	responses.WriteResponse(w, models.Webhook{
		ID:        newWebhookID,
		URL:       request.URL,
		Events:    request.Events,
		Timestamp: time.Now(),
		LastUsed:  nil,
	})
}

// GetWebhooks will return all webhooks.
func (h *Handlers) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	webhookRepository := storage.GetWebhookRepository()

	webhooks, err := webhookRepository.GetWebhooks()
	if err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	responses.WriteResponse(w, webhooks)
}

// DeleteWebhook will delete a single webhook.
func (h *Handlers) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request deleteWebhookRequest
	if err := decoder.Decode(&request); err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	webhookRepository := storage.GetWebhookRepository()

	if err := webhookRepository.DeleteWebhook(request.ID); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	responses.WriteSimpleResponse(w, true, "deleted webhook")
}
