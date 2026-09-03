package schedule

import (
	"errors"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
)

type federationScheduleRepo struct {
	scheduleeventsrepository.ScheduleEventsRepository
	events  []models.ScheduledEvent
	sentAt  map[string]time.Time
	loadErr error
}

func (r *federationScheduleRepo) GetEventsToFederate(time.Time) ([]models.ScheduledEvent, error) {
	if r.loadErr != nil {
		return nil, r.loadErr
	}
	var pending []models.ScheduledEvent
	for _, event := range r.events {
		if _, sent := r.sentAt[event.ID]; !sent {
			pending = append(pending, event)
		}
	}
	return pending, nil
}

func (r *federationScheduleRepo) SetEventFederatedAt(id string, sentAt time.Time) error {
	r.sentAt[id] = sentAt
	return nil
}

func TestPublishPendingEventsMarksOnlyQueuedEvents(t *testing.T) {
	now := time.Date(2030, time.September, 8, 12, 0, 0, 0, time.UTC)
	repo := &federationScheduleRepo{
		events: []models.ScheduledEvent{{ID: "first"}, {ID: "retry"}},
		sentAt: make(map[string]time.Time),
	}
	var sent []string
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return true },
		GetFederationEnabled:     func() bool { return true },
		FederateScheduledEvent: func(event models.ScheduledEvent) error {
			sent = append(sent, event.ID)
			if event.ID == "retry" {
				return errors.New("queue unavailable")
			}
			return nil
		},
	})

	service.publishPendingEvents(now)
	service.publishPendingEvents(now.Add(time.Minute))

	if len(sent) != 3 || sent[0] != "first" || sent[1] != "retry" || sent[2] != "retry" {
		t.Fatalf("sent = %#v", sent)
	}
	if got := repo.sentAt["first"]; !got.Equal(now) {
		t.Errorf("first federated at %v, want %v", got, now)
	}
	if _, stamped := repo.sentAt["retry"]; stamped {
		t.Error("failed event was marked federated")
	}
}

func TestPublishPendingEventsRequiresEnabledScheduleAndFederation(t *testing.T) {
	repo := &federationScheduleRepo{
		events: []models.ScheduledEvent{{ID: "event"}},
		sentAt: make(map[string]time.Time),
	}
	called := false
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return true },
		GetFederationEnabled:     func() bool { return false },
		FederateScheduledEvent: func(models.ScheduledEvent) error {
			called = true
			return nil
		},
	})

	service.publishPendingEvents(time.Now())
	if called {
		t.Error("disabled federation published a scheduled event")
	}
}
