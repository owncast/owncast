package webhooks

import "github.com/owncast/owncast/models"

type Manager interface {
	SendEventToWebhooks(payload WebhookEvent)
}

// LiveWebhookManager represents a central place for requesting webhooks
// to be sent out to all registered webhook destinations.
type LiveWebhookManager struct {
	queue     chan Job
	getStatus func() models.Status
}

// New creates a new webhook manager.
func New(getStatusFunc func() models.Status) *LiveWebhookManager {
	m := &LiveWebhookManager{
		getStatus: getStatusFunc,
	}
	m.initWorkerPool()
	return m
}

// InitTemporarySingleton initializes the the temporary global instance of the webhook manager
// to be deleted once dependency injection is implemented.
func InitTemporarySingleton(getStatusFunc func() models.Status) {
	temporaryGlobalInstance = New(getStatusFunc)
}

var temporaryGlobalInstance *LiveWebhookManager

// Get returns the temporary global instance of the webhook manager.
// Remove this after dependency injection is implemented.
func Get() *LiveWebhookManager {
	return temporaryGlobalInstance
}
