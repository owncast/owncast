package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/controllers"
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
func (c *Controller) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request createWebhookRequest
	if err := decoder.Decode(&request); err != nil {
		c.Service.BadRequestHandler(w, err)
		return
	}

	// Verify all the scopes provided are valid
	if !models.HasValidEvents(request.Events) {
		c.Service.BadRequestHandler(w, errors.New("one or more invalid event provided"))
		return
	}

	newWebhookID, err := c.Data.InsertWebhook(request.URL, request.Events)
	if err != nil {
		c.Service.InternalErrorHandler(w, err)
		return
	}

	c.Service.WriteResponse(w, models.Webhook{
		ID:        newWebhookID,
		URL:       request.URL,
		Events:    request.Events,
		Timestamp: time.Now(),
		LastUsed:  nil,
	})
}

// GetWebhooks will return all webhooks.
func (c *Controller) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := c.Data.GetWebhooks()
	if err != nil {
		c.Service.InternalErrorHandler(w, err)
		return
	}

	c.Service.WriteResponse(w, webhooks)
}

// DeleteWebhook will delete a single webhook.
func (c *Controller) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != controllers.POST {
		c.Service.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request deleteWebhookRequest
	if err := decoder.Decode(&request); err != nil {
		c.Service.BadRequestHandler(w, err)
		return
	}

	if err := c.Data.DeleteWebhook(request.ID); err != nil {
		c.Service.InternalErrorHandler(w, err)
		return
	}

	c.Service.WriteSimpleResponse(w, true, "deleted webhook")
}
