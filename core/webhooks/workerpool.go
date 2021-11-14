package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

const (
	// webhookWorkerPoolSize defines the number of concurrent HTTP webhook requests.
	webhookWorkerPoolSize = 10
)

// Job struct bundling the webhook and the payload in one struct.
type Job struct {
	webhook models.Webhook
	payload WebhookEvent
}

var queue chan Job

// InitWorkerPool starts n go routines that await webhook jobs.
func InitWorkerPool() {
	queue = make(chan Job)

	// start workers
	for i := 1; i <= webhookWorkerPoolSize; i++ {
		go worker(i, queue)
	}
}

func addToQueue(webhook models.Webhook, payload WebhookEvent) {
	log.Tracef("Queued Event %s for Webhook %s", payload.Type, webhook.URL)
	queue <- Job{webhook, payload}
}

func worker(workerID int, queue <-chan Job) {
	log.Debugf("Started Webhook worker %d", workerID)

	for job := range queue {
		log.Debugf("Event %s sent to Webhook %s using worker %d", job.payload.Type, job.webhook.URL, workerID)

		if err := sendWebhook(job); err != nil {
			log.Errorf("Event: %s failed to send to webhook: %s Error: %s", job.payload.Type, job.webhook.URL, err)
		}
		log.Tracef("Done with Event %s to Webhook %s using worker %d", job.payload.Type, job.webhook.URL, workerID)
	}
}

func sendWebhook(job Job) error {
	jsonText, err := json.Marshal(job.payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", job.webhook.URL, bytes.NewReader(jsonText))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := data.SetWebhookAsUsed(job.webhook); err != nil {
		log.Warnln(err)
	}

	return nil
}
