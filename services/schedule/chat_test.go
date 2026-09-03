package schedule

import (
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
)

type cancellingScheduleRepo struct {
	scheduleeventsrepository.ScheduleEventsRepository
	cancelled []string
}

func (r *cancellingScheduleRepo) CancelEvent(id string) error {
	r.cancelled = append(r.cancelled, id)
	return nil
}

func (r *cancellingScheduleRepo) GetCurrentOrUpcomingEvents(time.Time, int) ([]models.ScheduledEvent, error) {
	return nil, nil
}

func TestIsChatOpenForEvent(t *testing.T) {
	start := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	event := &models.ScheduledEvent{StartTime: start, DurationMinutes: 60}

	tests := []struct {
		name       string
		leadTime   time.Duration
		now        time.Time
		expectOpen bool
	}{
		{name: "disabled before start", leadTime: 0, now: start.Add(-time.Minute), expectOpen: false},
		{name: "opens at configured boundary", leadTime: 30 * time.Minute, now: start.Add(-30 * time.Minute), expectOpen: true},
		{name: "closed before configured boundary", leadTime: 30 * time.Minute, now: start.Add(-31 * time.Minute), expectOpen: false},
		{name: "open during event", leadTime: 60 * time.Minute, now: start.Add(30 * time.Minute), expectOpen: true},
		{name: "closed after event", leadTime: 60 * time.Minute, now: start.Add(61 * time.Minute), expectOpen: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if open := IsChatOpenForEvent(event, test.now, test.leadTime); open != test.expectOpen {
				t.Fatalf("expected chat open=%t, got %t", test.expectOpen, open)
			}
		})
	}
}

func TestIsEventMissed(t *testing.T) {
	start := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	event := &models.ScheduledEvent{StartTime: start}

	if IsEventMissed(event, start.Add(MissedEventGracePeriod-time.Nanosecond)) {
		t.Fatal("event reported missed before grace period elapsed")
	}
	if !IsEventMissed(event, start.Add(MissedEventGracePeriod)) {
		t.Fatal("event not reported missed at grace period boundary")
	}
	if IsEventMissedWarning(event, start.Add(MissedEventGracePeriod-MissedEventWarningLeadTime-time.Nanosecond)) {
		t.Fatal("event warning reported before warning boundary")
	}
	if !IsEventMissedWarning(event, start.Add(MissedEventGracePeriod-MissedEventWarningLeadTime)) {
		t.Fatal("event warning not reported one minute before shutdown")
	}
}

func TestUpdateChatWindowNotifiesOnceForMissedPreopenedEvent(t *testing.T) {
	start := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC)
	event := &models.ScheduledEvent{
		ID:              "event-1",
		StartTime:       start,
		DurationMinutes: 5,
	}
	missed := 0
	repo := &cancellingScheduleRepo{}
	service := New(Deps{
		ScheduleEventsRepository: repo,
		GetStatus: func() models.Status {
			return models.Status{}
		},
		GetChatOpenMinutes: func() int {
			return 10
		},
		OnMissedEventWarning: func(got *models.ScheduledEvent) {
			if got.ID == event.ID {
				missed++
			}
		},
	})
	service.upcomingEvent = event
	service.upcomingFetchedAt = time.Now()

	service.updateChatWindow(start.Add(-time.Minute))
	service.updateChatWindow(start.Add(MissedEventGracePeriod - MissedEventWarningLeadTime))
	service.updateChatWindow(start.Add(MissedEventGracePeriod))
	service.updateChatWindow(start.Add(MissedEventGracePeriod + time.Minute))

	if missed != 1 {
		t.Fatalf("missed event callback count = %d, want 1", missed)
	}
	if len(repo.cancelled) != 1 || repo.cancelled[0] != event.ID {
		t.Fatalf("cancelled events = %v, want [%s]", repo.cancelled, event.ID)
	}
}
