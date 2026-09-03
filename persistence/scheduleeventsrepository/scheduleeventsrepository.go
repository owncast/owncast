package scheduleeventsrepository

import (
	"time"

	"github.com/owncast/owncast/models"
)

const (
	ReminderFirst  = 1
	ReminderSecond = 2
)

// ScheduleEventsRepository defines the interface for scheduled stream event
// storage. All timestamps are normalized to UTC at this boundary so SQLite
// string comparisons on TIMESTAMP columns stay consistent.
type ScheduleEventsRepository interface {
	// Series (recurring schedules).
	AddSeries(name, description, reminderMessage, recurrence string, durationMinutes int) (string, error)
	GetSeries(id string) (*models.ScheduledEventSeries, error)
	GetAllSeries() ([]models.ScheduledEventSeries, error)
	GetActiveSeries() ([]models.ScheduledEventSeries, error)
	UpdateSeries(id, name, description, reminderMessage, recurrence string, durationMinutes int) error
	SetSeriesActive(id string, active bool) error
	DeleteSeries(id string) error

	// Occurrences (concrete events, one-off or materialized).
	AddOneOffEvent(name, description, reminderMessage string, start time.Time, durationMinutes int, timezone string) (string, error)
	// AddOccurrence inserts a materialized occurrence for a series. Returns
	// false when the (series, originalStart) slot already exists, making
	// repeat materialization a no-op.
	AddOccurrence(seriesID string, originalStart time.Time, name, description, reminderMessage string, start time.Time, durationMinutes int, timezone string) (bool, error)
	GetEvent(id string) (*models.ScheduledEvent, error)
	GetEventsInRange(from, to time.Time) ([]models.ScheduledEvent, error)
	GetEventsForSeries(seriesID string) ([]models.ScheduledEvent, error)
	UpdateEventDetails(id, name, description, reminderMessage string, durationMinutes int) error
	CancelEvent(id string) error
	// MoveEvent changes an occurrence's start time. Its original_start
	// identity is untouched so the materializer will not re-insert the
	// vacated slot.
	MoveEvent(id string, newStart time.Time) error
	DeleteEvent(id string) error
	// DeleteUnfederatedFutureEventsForSeries clears future occurrences nobody
	// has been told about, so an edited series can regenerate them. Announced
	// (federated) rows are kept for Update/Delete activities.
	DeleteUnfederatedFutureEventsForSeries(seriesID string, after time.Time) error
	// GetCurrentOrUpcomingEvents returns scheduled occurrences that are
	// still running or have not started, soonest first.
	GetCurrentOrUpcomingEvents(now time.Time, limit int) ([]models.ScheduledEvent, error)
	GetNextUpcomingEvents(after time.Time, limit int) ([]models.ScheduledEvent, error)
	GetEventsToFederate(startingAfter time.Time) ([]models.ScheduledEvent, error)
	SetEventFederatedAt(id string, t time.Time) error
	GetEventsNeedingReminder(startAfter, startBefore time.Time, reminderNumber int) ([]models.ScheduledEvent, error)
	SetEventReminderSentAt(id string, reminderNumber int, t time.Time) error
	GetEventsNeedingWebhookWarning(startAfter, startBefore time.Time) ([]models.ScheduledEvent, error)
	GetEventsNeedingWebhookStart(now time.Time) ([]models.ScheduledEvent, error)
	GetEventsNeedingWebhookEnd(now time.Time) ([]models.ScheduledEvent, error)
	SetEventWebhookWarningSentAt(id string, t time.Time) error
	SetEventWebhookStartedSentAt(id string, t time.Time) error
	SetEventWebhookEndedSentAt(id string, t time.Time) error
}
