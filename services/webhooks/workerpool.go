package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/dispatcher"
)

// Job bundles a single webhook destination + payload for a worker.
type Job struct {
	wg      *sync.WaitGroup
	payload WebhookEvent
	webhook models.Webhook
}

// Service owns the webhook dispatcher: a bounded worker pool that
// fans events out to configured destinations.
//
// Most events also need the current stream status and follower details
// to build their payloads, so those collaborators are dependencies.
type Service struct {
	workerPoolSize int
	queue          chan Job

	// getStatus returns the current stream status for inclusion in
	// outbound event payloads. Supplied by main.go (the stream service).
	getStatus func() models.Status

	// followers is consulted by the fediverse-engagement builder to
	// resolve actor IRIs to display data.
	followers followersrepository.FollowersRepository

	// configRepository provides server metadata included in stream-status
	// event payloads.
	configRepository configrepository.ConfigRepository

	// webhookRepository stores configured webhook destinations and is
	// consulted when dispatching events as well as when marking a
	// destination as used after a successful send.
	webhookRepository webhookrepository.WebhookRepository

	// events is the shared dispatcher every event is published to, so
	// in-process consumers (the plugin host today, other subsystems later)
	// receive them without this service knowing who they are.
	events *dispatcher.Dispatcher
}

// Deps lists every collaborator a *Service needs at construction.
type Deps struct {
	GetStatus         func() models.Status
	Followers         followersrepository.FollowersRepository
	ConfigRepository  configrepository.ConfigRepository
	WebhookRepository webhookrepository.WebhookRepository
	Events            *dispatcher.Dispatcher
}

// New constructs an idle webhook Service. Call Start to launch the
// worker pool.
func New(deps Deps) *Service {
	return &Service{
		workerPoolSize:    runtime.GOMAXPROCS(0),
		getStatus:         deps.GetStatus,
		followers:         deps.Followers,
		configRepository:  deps.ConfigRepository,
		webhookRepository: deps.WebhookRepository,
		events:            deps.Events,
	}
}

// SetGetStatus wires the stream-status callback after construction.
// Exists because the composition root has a small cycle: webhooks
// needs stream.GetStatus, but stream needs *webhooks.Service. main.go
// constructs webhooks first with a nil callback, then fills it in once
// streamSvc exists. Must be called before Start.
func (s *Service) SetGetStatus(fn func() models.Status) {
	s.getStatus = fn
}

// Start launches the worker goroutines that drain the queue.
func (s *Service) Start() {
	s.queue = make(chan Job)
	for i := 1; i <= s.workerPoolSize; i++ {
		go s.worker(i)
	}
}

func (s *Service) addToQueue(webhook models.Webhook, payload WebhookEvent, wg *sync.WaitGroup) {
	log.Tracef("Queued Event %s for Webhook %s", payload.Type, webhook.URL)
	s.queue <- Job{wg, payload, webhook}
}

func (s *Service) worker(workerID int) {
	log.Debugf("Started Webhook worker %d", workerID)

	for job := range s.queue {
		log.Debugf("Event %s sent to Webhook %s using worker %d", job.payload.Type, job.webhook.URL, workerID)

		if err := s.sendWebhook(job); err != nil {
			log.Errorf("Event: %s failed to send to webhook: %s Error: %s", job.payload.Type, job.webhook.URL, err)
		}
		log.Tracef("Done with Event %s to Webhook %s using worker %d", job.payload.Type, job.webhook.URL, workerID)
		if job.wg != nil {
			job.wg.Done()
		}
	}
}

func (s *Service) sendWebhook(job Job) error {
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

	if err := s.webhookRepository.SetWebhookAsUsed(job.webhook); err != nil {
		log.Warnln(err)
	}

	return nil
}
