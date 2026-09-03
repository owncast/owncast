package schedule

import (
	"context"
	"errors"
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

func (lifecycleWebhookConfig) GetServerURL() string               { return "https://owncast.example" }
func (lifecycleWebhookConfig) GetScheduleReminderMessage() string { return "" }

func TestScheduledEventWebhooksFireOnce(t *testing.T) {
	now := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	repo := &lifecycleScheduleRepo{
		warning: models.ScheduledEvent{ID: "warning"},
		started: models.ScheduledEvent{ID: "started", Name: "Scheduled title"},
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
	var titles []string
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return true },
		SetStreamTitle: func(title string) error {
			titles = append(titles, title)
			return nil
		},
		Webhooks: webhookService,
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
	if len(titles) != 1 || titles[0] != "Scheduled title" {
		t.Errorf("stream titles = %v, want [Scheduled title]", titles)
	}
}

func TestScheduledEventTitleFailureRetriesWithoutStartingLifecycle(t *testing.T) {
	now := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	repo := &lifecycleScheduleRepo{started: models.ScheduledEvent{ID: "started", Name: "Scheduled title"}}
	events := dispatcher.New()
	var actions []models.ScheduledEventWebhookAction
	events.AddListener(func(_ context.Context, event dispatcher.Event) {
		data := event.Payload.(webhooks.WebhookEvent).EventData.(*webhooks.WebhookScheduledEventData)
		actions = append(actions, data.Action)
	})
	webhookService := webhooks.New(webhooks.Deps{
		GetStatus:         func() models.Status { return models.Status{} },
		ConfigRepository:  lifecycleWebhookConfig{},
		WebhookRepository: lifecycleWebhookRepo{},
		Events:            events,
	})
	attempts := 0
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return true },
		SetStreamTitle: func(string) error {
			attempts++
			if attempts == 1 {
				return errors.New("storage unavailable")
			}
			return nil
		},
		Webhooks: webhookService,
	})

	service.sendScheduledEventWebhooks(now)
	if len(actions) != 0 {
		t.Fatalf("actions after title failure = %v, want none", actions)
	}
	service.sendScheduledEventWebhooks(now.Add(time.Minute))
	if len(actions) != 1 || actions[0] != models.ScheduledEventStarted {
		t.Fatalf("actions after title retry = %v, want [started]", actions)
	}
	if attempts != 2 {
		t.Errorf("title attempts = %d, want 2", attempts)
	}
}
