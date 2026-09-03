package admin

import (
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	"github.com/owncast/owncast/services/schedule"
)

type federationMutationRepo struct {
	scheduleeventsrepository.ScheduleEventsRepository
	events         map[string]models.ScheduledEvent
	deleted        []string
	pendingDeletes []string
}

func (r *federationMutationRepo) UpdateEventFromSeries(id, name, description, reminderMessage string, durationMinutes int, timezone string) error {
	event := r.events[id]
	event.Name = name
	event.Description = description
	event.ReminderMessage = reminderMessage
	event.DurationMinutes = durationMinutes
	event.Timezone = timezone
	updatedAt := time.Date(2030, time.September, 9, 0, 0, 0, 0, time.UTC)
	event.UpdatedAt = &updatedAt
	r.events[id] = event
	return nil
}

func (r *federationMutationRepo) GetEvent(id string) (*models.ScheduledEvent, error) {
	event, ok := r.events[id]
	if !ok {
		return nil, nil
	}
	return &event, nil
}

func (r *federationMutationRepo) DeleteEvent(id string) error {
	if event, ok := r.events[id]; ok && event.FederatedAt != nil {
		r.pendingDeletes = append(r.pendingDeletes, id)
	}
	delete(r.events, id)
	r.deleted = append(r.deleted, id)
	return nil
}

func (r *federationMutationRepo) ClearEventFederationUpdatePending(string, int64) error {
	return nil
}

func (r *federationMutationRepo) GetPendingFederationDeletes() ([]string, error) {
	return r.pendingDeletes, nil
}

func (r *federationMutationRepo) ClearPendingFederationDelete(id string) error {
	for i, pending := range r.pendingDeletes {
		if pending == id {
			r.pendingDeletes = append(r.pendingDeletes[:i], r.pendingDeletes[i+1:]...)
			break
		}
	}
	return nil
}

func TestSyncFederatedSeriesEventsUpdatesRetainedAndDeletesRemoved(t *testing.T) {
	now := time.Date(2030, time.September, 8, 12, 0, 0, 0, time.UTC)
	federatedAt := now.Add(-time.Hour)
	retainedStart := time.Date(2030, time.September, 10, 12, 0, 0, 0, time.UTC)
	removedStart := retainedStart.Add(24 * time.Hour)
	retainedOriginal := retainedStart
	removedOriginal := removedStart
	retained := models.ScheduledEvent{ID: "retained", OriginalStart: &retainedOriginal, StartTime: retainedStart, Status: models.ScheduledEventStatusScheduled, FederatedAt: &federatedAt}
	removed := models.ScheduledEvent{ID: "removed", OriginalStart: &removedOriginal, StartTime: removedStart, Status: models.ScheduledEventStatusScheduled, FederatedAt: &federatedAt}
	repo := &federationMutationRepo{events: map[string]models.ScheduledEvent{retained.ID: retained, removed.ID: removed}}
	var updated, deleted []string
	scheduleService := schedule.New(schedule.Deps{
		ScheduleEventsRepository: repo,
		GetScheduleEnabled:       func() bool { return true },
		GetFederationEnabled:     func() bool { return true },
		FederateScheduledEventUpdate: func(event models.ScheduledEvent) error {
			updated = append(updated, event.ID)
			return nil
		},
		FederateScheduledEventDelete: func(event models.ScheduledEvent) error {
			deleted = append(deleted, event.ID)
			return nil
		},
	})
	admin := &Admin{scheduleEventsRepository: repo, schedule: scheduleService}
	series := models.ScheduledEventSeries{
		ID:              "series",
		Name:            "Updated series",
		Description:     "Updated description",
		ReminderMessage: "Soon",
		Recurrence:      "DTSTART;TZID=UTC:20300910T120000\nRRULE:FREQ=DAILY;COUNT=1",
		DurationMinutes: 90,
	}

	if err := admin.syncFederatedSeriesEvents([]models.ScheduledEvent{retained, removed}, series, now); err != nil {
		t.Fatal(err)
	}
	if len(updated) != 1 || updated[0] != retained.ID {
		t.Fatalf("updated = %#v", updated)
	}
	if len(deleted) != 1 || deleted[0] != removed.ID || len(repo.deleted) != 1 || repo.deleted[0] != removed.ID {
		t.Fatalf("deleted activities/rows = %#v / %#v", deleted, repo.deleted)
	}
	current := repo.events[retained.ID]
	if current.Name != series.Name || current.DurationMinutes != series.DurationMinutes || current.Timezone != "UTC" {
		t.Fatalf("retained event = %#v", current)
	}
}
