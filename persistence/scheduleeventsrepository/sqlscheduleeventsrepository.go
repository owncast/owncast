package scheduleeventsrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/datastore"
	"github.com/teris-io/shortid"
)

// SqlScheduleEventsRepository is a SQL implementation of the
// ScheduleEventsRepository interface.
type SqlScheduleEventsRepository struct {
	datastore *datastore.Datastore
}

// temporaryGlobalInstance is set once during application startup so helper
// code that has not yet been migrated to the dependency-injection pattern
// can still reach the schedule-events repository.
var temporaryGlobalInstance ScheduleEventsRepository

// SetGlobalInstance registers the application's single
// ScheduleEventsRepository for Get to return. Called from main.go after
// constructing the repository.
func SetGlobalInstance(r ScheduleEventsRepository) {
	temporaryGlobalInstance = r
}

// Get returns the global ScheduleEventsRepository registered with
// SetGlobalInstance. Returns nil until startup has wired one in.
func Get() ScheduleEventsRepository {
	return temporaryGlobalInstance
}

// New will create a new instance of the ScheduleEventsRepository.
func New(datastore *datastore.Datastore) ScheduleEventsRepository {
	return &SqlScheduleEventsRepository{
		datastore: datastore,
	}
}

// AddSeries creates a recurring series and returns its generated id.
func (r *SqlScheduleEventsRepository) AddSeries(name, description, recurrence string, durationMinutes int) (string, error) {
	id := shortid.MustGenerate()
	queries := db.New(r.datastore.DB)
	err := queries.AddStreamEventSeries(context.Background(), db.AddStreamEventSeriesParams{
		ID:              id,
		Name:            name,
		Description:     description,
		Recurrence:      recurrence,
		DurationMinutes: int64(durationMinutes),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetSeries returns one series by id, or nil when it does not exist.
func (r *SqlScheduleEventsRepository) GetSeries(id string) (*models.ScheduledEventSeries, error) {
	queries := db.New(r.datastore.DB)
	row, err := queries.GetStreamEventSeries(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	series := seriesFromRow(row)
	return &series, nil
}

// GetAllSeries returns every series, recurring schedules paused or not.
func (r *SqlScheduleEventsRepository) GetAllSeries() ([]models.ScheduledEventSeries, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetAllStreamEventSeries(context.Background())
	if err != nil {
		return nil, err
	}
	series := make([]models.ScheduledEventSeries, 0, len(rows))
	for _, row := range rows {
		series = append(series, seriesFromRow(row))
	}
	return series, nil
}

// GetActiveSeries returns the series the materializer should expand.
func (r *SqlScheduleEventsRepository) GetActiveSeries() ([]models.ScheduledEventSeries, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetActiveStreamEventSeries(context.Background())
	if err != nil {
		return nil, err
	}
	series := make([]models.ScheduledEventSeries, 0, len(rows))
	for _, row := range rows {
		series = append(series, seriesFromRow(row))
	}
	return series, nil
}

// UpdateSeries rewrites a series' definition. The caller is responsible for
// regenerating future occurrences afterwards.
func (r *SqlScheduleEventsRepository) UpdateSeries(id, name, description, recurrence string, durationMinutes int) error {
	queries := db.New(r.datastore.DB)
	return queries.UpdateStreamEventSeries(context.Background(), db.UpdateStreamEventSeriesParams{
		Name:            name,
		Description:     description,
		Recurrence:      recurrence,
		DurationMinutes: int64(durationMinutes),
		ID:              id,
	})
}

// SetSeriesActive pauses or resumes materialization for a series.
func (r *SqlScheduleEventsRepository) SetSeriesActive(id string, active bool) error {
	queries := db.New(r.datastore.DB)
	return queries.SetStreamEventSeriesActive(context.Background(), db.SetStreamEventSeriesActiveParams{
		Active: active,
		ID:     id,
	})
}

// DeleteSeries removes a series row. Its occurrences are the caller's
// responsibility (unfederated ones get deleted, announced ones get Delete
// activities).
func (r *SqlScheduleEventsRepository) DeleteSeries(id string) error {
	queries := db.New(r.datastore.DB)
	return queries.DeleteStreamEventSeries(context.Background(), id)
}

// AddOneOffEvent creates a standalone occurrence with no series and returns
// its generated id.
func (r *SqlScheduleEventsRepository) AddOneOffEvent(name, description string, start time.Time, durationMinutes int, timezone string) (string, error) {
	id := shortid.MustGenerate()
	queries := db.New(r.datastore.DB)
	inserted, err := queries.AddStreamEvent(context.Background(), db.AddStreamEventParams{
		ID:              id,
		SeriesID:        sql.NullString{},
		OriginalStart:   sql.NullTime{},
		Name:            name,
		Description:     description,
		StartTime:       start.UTC(),
		DurationMinutes: int64(durationMinutes),
		Timezone:        timezone,
	})
	if err != nil {
		return "", err
	}
	// A one-off has no legitimate INSERT OR IGNORE path (NULL series and
	// original_start never collide in the unique index), so an ignored
	// insert can only mean the generated id hit an existing row. Returning
	// that id would point the caller at somebody else's event.
	if inserted == 0 {
		return "", errors.New("scheduled event id collision")
	}
	return id, nil
}

// AddOccurrence inserts a materialized occurrence for a series. Returns
// false when the (series, originalStart) slot already exists.
func (r *SqlScheduleEventsRepository) AddOccurrence(seriesID string, originalStart time.Time, name, description string, start time.Time, durationMinutes int, timezone string) (bool, error) {
	id := shortid.MustGenerate()
	queries := db.New(r.datastore.DB)
	inserted, err := queries.AddStreamEvent(context.Background(), db.AddStreamEventParams{
		ID:              id,
		SeriesID:        sql.NullString{String: seriesID, Valid: true},
		OriginalStart:   models.TimeToNullTime(originalStart.UTC()),
		Name:            name,
		Description:     description,
		StartTime:       start.UTC(),
		DurationMinutes: int64(durationMinutes),
		Timezone:        timezone,
	})
	if err != nil {
		return false, err
	}
	return inserted > 0, nil
}

// GetEvent returns one occurrence by id, or nil when it does not exist.
func (r *SqlScheduleEventsRepository) GetEvent(id string) (*models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	row, err := queries.GetStreamEvent(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	event := eventFromRow(row)
	return &event, nil
}

// GetEventsInRange returns occurrences with start_time in [from, to).
func (r *SqlScheduleEventsRepository) GetEventsInRange(from, to time.Time) ([]models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetStreamEventsInRange(context.Background(), db.GetStreamEventsInRangeParams{
		StartTime:   from.UTC(),
		StartTime_2: to.UTC(),
	})
	if err != nil {
		return nil, err
	}
	return eventsFromRows(rows), nil
}

// GetEventsForSeries returns every occurrence belonging to a series.
func (r *SqlScheduleEventsRepository) GetEventsForSeries(seriesID string) ([]models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetStreamEventsForSeries(context.Background(), sql.NullString{String: seriesID, Valid: true})
	if err != nil {
		return nil, err
	}
	return eventsFromRows(rows), nil
}

// UpdateEventDetails rewrites an occurrence's descriptive fields.
func (r *SqlScheduleEventsRepository) UpdateEventDetails(id, name, description string, durationMinutes int) error {
	queries := db.New(r.datastore.DB)
	return queries.UpdateStreamEventDetails(context.Background(), db.UpdateStreamEventDetailsParams{
		Name:            name,
		Description:     description,
		DurationMinutes: int64(durationMinutes),
		ID:              id,
	})
}

// CancelEvent marks an occurrence cancelled, keeping the row.
func (r *SqlScheduleEventsRepository) CancelEvent(id string) error {
	queries := db.New(r.datastore.DB)
	return queries.CancelStreamEvent(context.Background(), id)
}

// MoveEvent changes an occurrence's start time, keeping its original_start
// identity.
func (r *SqlScheduleEventsRepository) MoveEvent(id string, newStart time.Time) error {
	queries := db.New(r.datastore.DB)
	return queries.MoveStreamEvent(context.Background(), db.MoveStreamEventParams{
		StartTime: newStart.UTC(),
		ID:        id,
	})
}

// DeleteEvent removes an occurrence row entirely.
func (r *SqlScheduleEventsRepository) DeleteEvent(id string) error {
	queries := db.New(r.datastore.DB)
	return queries.DeleteStreamEvent(context.Background(), id)
}

// DeleteUnfederatedFutureEventsForSeries clears regenerable occurrences after
// a series edit.
func (r *SqlScheduleEventsRepository) DeleteUnfederatedFutureEventsForSeries(seriesID string, after time.Time) error {
	queries := db.New(r.datastore.DB)
	return queries.DeleteUnfederatedFutureStreamEventsForSeries(context.Background(), db.DeleteUnfederatedFutureStreamEventsForSeriesParams{
		SeriesID:  sql.NullString{String: seriesID, Valid: true},
		StartTime: after.UTC(),
	})
}

// GetCurrentOrUpcomingEvents returns scheduled occurrences that are still
// running (start plus duration in the future) or upcoming, soonest first.
func (r *SqlScheduleEventsRepository) GetCurrentOrUpcomingEvents(now time.Time, limit int) ([]models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetCurrentOrUpcomingStreamEvents(context.Background(), db.GetCurrentOrUpcomingStreamEventsParams{
		Datetime: now.UTC(),
		Limit:    int64(limit),
	})
	if err != nil {
		return nil, err
	}
	return eventsFromRows(rows), nil
}

// GetNextUpcomingEvents returns the next scheduled (not cancelled)
// occurrences starting after the given instant.
func (r *SqlScheduleEventsRepository) GetNextUpcomingEvents(after time.Time, limit int) ([]models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetNextUpcomingStreamEvents(context.Background(), db.GetNextUpcomingStreamEventsParams{
		StartTime: after.UTC(),
		Limit:     int64(limit),
	})
	if err != nil {
		return nil, err
	}
	return eventsFromRows(rows), nil
}

// GetEventsToFederate returns future scheduled occurrences whose Create has
// not gone out yet.
func (r *SqlScheduleEventsRepository) GetEventsToFederate(startingAfter time.Time) ([]models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetStreamEventsToFederate(context.Background(), startingAfter.UTC())
	if err != nil {
		return nil, err
	}
	return eventsFromRows(rows), nil
}

// SetEventFederatedAt stamps an occurrence as announced.
func (r *SqlScheduleEventsRepository) SetEventFederatedAt(id string, t time.Time) error {
	queries := db.New(r.datastore.DB)
	return queries.SetStreamEventFederatedAt(context.Background(), db.SetStreamEventFederatedAtParams{
		FederatedAt: models.TimeToNullTime(t.UTC()),
		ID:          id,
	})
}

// GetEventsNeedingReminder returns unreminded scheduled occurrences starting
// inside (startAfter, startBefore].
func (r *SqlScheduleEventsRepository) GetEventsNeedingReminder(startAfter, startBefore time.Time) ([]models.ScheduledEvent, error) {
	queries := db.New(r.datastore.DB)
	rows, err := queries.GetStreamEventsNeedingReminder(context.Background(), db.GetStreamEventsNeedingReminderParams{
		StartTime:   startAfter.UTC(),
		StartTime_2: startBefore.UTC(),
	})
	if err != nil {
		return nil, err
	}
	return eventsFromRows(rows), nil
}

// SetEventReminderSentAt stamps an occurrence's reminder as sent.
func (r *SqlScheduleEventsRepository) SetEventReminderSentAt(id string, t time.Time) error {
	queries := db.New(r.datastore.DB)
	return queries.SetStreamEventReminderSentAt(context.Background(), db.SetStreamEventReminderSentAtParams{
		ReminderSentAt: models.TimeToNullTime(t.UTC()),
		ID:             id,
	})
}

func seriesFromRow(row db.StreamEventSeries) models.ScheduledEventSeries {
	return models.ScheduledEventSeries{
		ID:              row.ID,
		Name:            row.Name,
		Description:     row.Description,
		Recurrence:      row.Recurrence,
		DurationMinutes: int(row.DurationMinutes),
		Active:          row.Active,
	}
}

func eventFromRow(row db.StreamEvent) models.ScheduledEvent {
	event := models.ScheduledEvent{
		ID:              row.ID,
		Name:            row.Name,
		Description:     row.Description,
		StartTime:       row.StartTime.UTC(),
		DurationMinutes: int(row.DurationMinutes),
		Timezone:        row.Timezone,
		Status:          row.Status,
	}
	if row.SeriesID.Valid {
		seriesID := row.SeriesID.String
		event.SeriesID = &seriesID
	}
	if row.OriginalStart.Valid {
		originalStart := row.OriginalStart.Time.UTC()
		event.OriginalStart = &originalStart
	}
	if row.FederatedAt.Valid {
		federatedAt := row.FederatedAt.Time.UTC()
		event.FederatedAt = &federatedAt
	}
	if row.ReminderSentAt.Valid {
		reminderSentAt := row.ReminderSentAt.Time.UTC()
		event.ReminderSentAt = &reminderSentAt
	}
	return event
}

func eventsFromRows(rows []db.StreamEvent) []models.ScheduledEvent {
	events := make([]models.ScheduledEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, eventFromRow(row))
	}
	return events
}
