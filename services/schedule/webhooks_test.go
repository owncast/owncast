package schedule

import (
	"context"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/webhooks"
)

type lifecycleScheduleRepo struct {
	scheduleeventsrepository.ScheduleEventsRepository
	warning models.ScheduledEvent
	started models.ScheduledEvent
	ended   models.ScheduledEvent
}

func (r *lifecycleScheduleRepo) GetEventsNeedingWebhookWarning(time.Time, time.Time) ([]models.ScheduledEvent, error) {
	if r.warning.ID == "" {
		return nil, nil
	}
	return []models.ScheduledEvent{r.warning}, nil
}

func (r *lifecycleScheduleRepo) GetEventsNeedingWebhookStart(time.Time) ([]models.ScheduledEvent, error) {
	if r.started.ID == "" {
		return nil, nil
	}
	return []models.ScheduledEvent{r.started}, nil
}

func (r *lifecycleScheduleRepo) GetEventsNeedingWebhookEnd(time.Time) ([]models.ScheduledEvent, error) {
	if r.ended.ID == "" {
		return nil, nil
	}
	return []models.ScheduledEvent{r.ended}, nil
}

func (r *lifecycleScheduleRepo) SetEventWebhookWarningSentAt(string, time.Time) error {
	r.warning = models.ScheduledEvent{}
	return nil
}

func (r *lifecycleScheduleRepo) SetEventWebhookStartedSentAt(string, time.Time) error {
	r.started = models.ScheduledEvent{}
	return nil
}

func (r *lifecycleScheduleRepo) SetEventWebhookEndedSentAt(string, time.Time) error {
	r.ended = models.ScheduledEvent{}
	return nil
}

type lifecycleWebhookRepo struct {
	webhookrepository.WebhookRepository
}

func (lifecycleWebhookRepo) GetWebhooksForEvent(models.EventType) []models.Webhook { return nil }

type lifecycleWebhookConfig struct {
	configrepository.ConfigRepository
}

func (lifecycleWebhookConfig) GetServerURL() string { return "https://owncast.example" }

func TestScheduledEventWebhooksFireOnce(t *testing.T) {
	now := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	repo := &lifecycleScheduleRepo{
		warning: models.ScheduledEvent{ID: "warning"},
		started: models.ScheduledEvent{ID: "started"},
		ended:   models.ScheduledEvent{ID: "ended"},
	}
	events := dispatcher.New()
	var actions []models.ScheduledEventWebhookAction
	events.AddListener(func(_ context.Context, event dispatcher.Event) {
		payload := event.Payload.(webhooks.WebhookEvent)
		data := payload.EventData.(*webhooks.WebhookScheduledEventData)
		actions = append(actions, data.Action)
	})
	webhookService := webhooks.New(webhooks.Deps{
		GetStatus:         func() models.Status { return models.Status{} },
		ConfigRepository:  lifecycleWebhookConfig{},
		WebhookRepository: lifecycleWebhookRepo{},
		Events:            events,
	})
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return true },
		Webhooks:                 webhookService,
	})

	service.sendScheduledEventWebhooks(now)
	service.sendScheduledEventWebhooks(now.Add(time.Minute))

	want := []models.ScheduledEventWebhookAction{
		models.ScheduledEventWarning,
		models.ScheduledEventStarted,
		models.ScheduledEventEnded,
	}
	if len(actions) != len(want) {
		t.Fatalf("webhook actions = %v, want %v", actions, want)
	}
	for i := range want {
		if actions[i] != want[i] {
			t.Errorf("webhook action %d = %q, want %q", i, actions[i], want[i])
		}
	}
}
