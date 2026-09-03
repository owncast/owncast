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
	events         []models.ScheduledEvent
	pendingUpdates []models.ScheduledEvent
	pendingDeletes []string
	sentAt         map[string]time.Time
	clearedUpdates []string
	clearedDeletes []string
	loadErr        error
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

func (r *federationScheduleRepo) SetEventFederatedAt(id string, sentAt time.Time, _ int64) error {
	r.sentAt[id] = sentAt
	return nil
}

func (r *federationScheduleRepo) GetEventsNeedingFederationUpdate() ([]models.ScheduledEvent, error) {
	return r.pendingUpdates, nil
}

func (r *federationScheduleRepo) ClearEventFederationUpdatePending(id string, _ int64) error {
	r.clearedUpdates = append(r.clearedUpdates, id)
	return nil
}

func (r *federationScheduleRepo) GetPendingFederationDeletes() ([]string, error) {
	return r.pendingDeletes, nil
}

func (r *federationScheduleRepo) ClearPendingFederationDelete(id string) error {
	r.clearedDeletes = append(r.clearedDeletes, id)
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

func TestPendingEventUpdatesRetryWhenScheduleIsDisabled(t *testing.T) {
	repo := &federationScheduleRepo{
		pendingUpdates: []models.ScheduledEvent{{ID: "event"}},
		pendingDeletes: []string{"event"},
	}
	attempts := 0
	deletes := 0
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return false },
		GetFederationEnabled:     func() bool { return true },
		FederateScheduledEventUpdate: func(models.ScheduledEvent) error {
			attempts++
			if attempts == 1 {
				return errors.New("queue unavailable")
			}
			return nil
		},
		FederateScheduledEventDelete: func(models.ScheduledEvent) error {
			deletes++
			return nil
		},
	})

	service.publishPendingEventUpdates()
	service.publishPendingEventUpdates()
	service.publishPendingEventDeletes()

	if attempts != 2 {
		t.Fatalf("update attempts = %d, want retry after queue failure", attempts)
	}
	if len(repo.clearedUpdates) != 1 || repo.clearedUpdates[0] != "event" {
		t.Fatalf("cleared updates = %#v, want event after successful queue", repo.clearedUpdates)
	}
	if deletes != 1 {
		t.Fatal("disabled schedule suppressed an Event deletion")
	}
}

func TestPendingEventDeleteWaitsForFederation(t *testing.T) {
	repo := &federationScheduleRepo{pendingDeletes: []string{"event"}}
	enabled := false
	deletes := 0
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetFederationEnabled:     func() bool { return enabled },
		FederateScheduledEventDelete: func(models.ScheduledEvent) error {
			deletes++
			return nil
		},
	})

	service.publishPendingEventDeletes()
	enabled = true
	service.publishPendingEventDeletes()

	if deletes != 1 {
		t.Fatalf("delete attempts = %d, want one after federation is enabled", deletes)
	}
	if len(repo.clearedDeletes) != 1 || repo.clearedDeletes[0] != "event" {
		t.Fatalf("cleared deletes = %#v", repo.clearedDeletes)
	}
}
