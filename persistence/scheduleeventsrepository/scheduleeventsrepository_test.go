package scheduleeventsrepository

import (
	"os"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/migrations"
	"github.com/owncast/owncast/services/datastore"
)

var testRepo ScheduleEventsRepository

func TestMain(m *testing.M) {
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		panic(err)
	}
	if err := migrations.Run(ds.DB, os.TempDir()); err != nil {
		panic(err)
	}

	configRepo := configrepository.New(ds)
	configRepo.SetServerURL("https://test.owncast.server")
	configRepo.SetFederationUsername("testuser")
	configrepository.SetGlobalInstance(configRepo)

	testRepo = New(ds)
	SetGlobalInstance(testRepo)

	m.Run()
}

func containsSeriesID(series []models.ScheduledEventSeries, id string) bool {
	for _, s := range series {
		if s.ID == id {
			return true
		}
	}
	return false
}

func containsEventID(events []models.ScheduledEvent, id string) bool {
	for _, e := range events {
		if e.ID == id {
			return true
		}
	}
	return false
}

func mustAddOneOffEvent(t *testing.T, name string, start time.Time) string {
	t.Helper()
	id, err := testRepo.AddOneOffEvent(name, "description for "+name, start, 60, "UTC")
	if err != nil {
		t.Fatalf("AddOneOffEvent(%s) unexpected error = %v", name, err)
	}
	return id
}

func TestSeriesCRUDRoundtrip(t *testing.T) {
	repo := Get()

	recurrence := "DTSTART;TZID=America/Los_Angeles:20270301T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO"
	id, err := repo.AddSeries("series-crud original", "first description", recurrence, 90)
	if err != nil {
		t.Fatalf("AddSeries() unexpected error = %v", err)
	}
	if id == "" {
		t.Fatalf("AddSeries() returned an empty id")
	}

	series, err := repo.GetSeries(id)
	if err != nil {
		t.Fatalf("GetSeries() unexpected error = %v", err)
	}
	if series == nil {
		t.Fatalf("GetSeries() returned nil for a just-added series")
	}
	if series.ID != id {
		t.Errorf("series ID = %v, want %v", series.ID, id)
	}
	if series.Name != "series-crud original" {
		t.Errorf("series Name = %v, want %v", series.Name, "series-crud original")
	}
	if series.Description != "first description" {
		t.Errorf("series Description = %v, want %v", series.Description, "first description")
	}
	if series.Recurrence != recurrence {
		t.Errorf("series Recurrence = %v, want %v", series.Recurrence, recurrence)
	}
	if series.DurationMinutes != 90 {
		t.Errorf("series DurationMinutes = %v, want 90", series.DurationMinutes)
	}

	allSeries, err := repo.GetAllSeries()
	if err != nil {
		t.Fatalf("GetAllSeries() unexpected error = %v", err)
	}
	if !containsSeriesID(allSeries, id) {
		t.Errorf("GetAllSeries() does not contain the added series %v", id)
	}

	updatedRecurrence := "DTSTART;TZID=America/New_York:20270302T200000\nRRULE:FREQ=WEEKLY;BYDAY=TU"
	if err := repo.UpdateSeries(id, "series-crud renamed", "second description", updatedRecurrence, 45); err != nil {
		t.Fatalf("UpdateSeries() unexpected error = %v", err)
	}
	series, err = repo.GetSeries(id)
	if err != nil || series == nil {
		t.Fatalf("GetSeries() after update = %v, %v", series, err)
	}
	if series.Name != "series-crud renamed" {
		t.Errorf("updated series Name = %v, want %v", series.Name, "series-crud renamed")
	}
	if series.Description != "second description" {
		t.Errorf("updated series Description = %v, want %v", series.Description, "second description")
	}
	if series.Recurrence != updatedRecurrence {
		t.Errorf("updated series Recurrence = %v, want %v", series.Recurrence, updatedRecurrence)
	}
	if series.DurationMinutes != 45 {
		t.Errorf("updated series DurationMinutes = %v, want 45", series.DurationMinutes)
	}

	// Active lifecycle: pause removes the series from the materializer's
	// view, resume brings it back.
	if err := repo.SetSeriesActive(id, false); err != nil {
		t.Fatalf("SetSeriesActive(false) unexpected error = %v", err)
	}
	series, err = repo.GetSeries(id)
	if err != nil || series == nil {
		t.Fatalf("GetSeries() after deactivate = %v, %v", series, err)
	}
	if series.Active {
		t.Errorf("series Active = true after SetSeriesActive(false)")
	}
	activeSeries, err := repo.GetActiveSeries()
	if err != nil {
		t.Fatalf("GetActiveSeries() unexpected error = %v", err)
	}
	if containsSeriesID(activeSeries, id) {
		t.Errorf("GetActiveSeries() contains a paused series %v", id)
	}

	if err := repo.SetSeriesActive(id, true); err != nil {
		t.Fatalf("SetSeriesActive(true) unexpected error = %v", err)
	}
	series, err = repo.GetSeries(id)
	if err != nil || series == nil {
		t.Fatalf("GetSeries() after reactivate = %v, %v", series, err)
	}
	if !series.Active {
		t.Errorf("series Active = false after SetSeriesActive(true)")
	}
	activeSeries, err = repo.GetActiveSeries()
	if err != nil {
		t.Fatalf("GetActiveSeries() unexpected error = %v", err)
	}
	if !containsSeriesID(activeSeries, id) {
		t.Errorf("GetActiveSeries() does not contain a resumed series %v", id)
	}

	if err := repo.DeleteSeries(id); err != nil {
		t.Fatalf("DeleteSeries() unexpected error = %v", err)
	}
	series, err = repo.GetSeries(id)
	if err != nil {
		t.Fatalf("GetSeries() after delete unexpected error = %v", err)
	}
	if series != nil {
		t.Errorf("GetSeries() after delete = %+v, want nil", series)
	}
	allSeries, err = repo.GetAllSeries()
	if err != nil {
		t.Fatalf("GetAllSeries() unexpected error = %v", err)
	}
	if containsSeriesID(allSeries, id) {
		t.Errorf("GetAllSeries() still contains the deleted series %v", id)
	}
}

func TestOneOffEventCRUDRoundtrip(t *testing.T) {
	repo := Get()

	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("LoadLocation() unexpected error = %v", err)
	}
	// A non-UTC input proves the instant survives the repository's UTC
	// normalization boundary.
	start := time.Date(2030, 8, 15, 11, 30, 0, 0, la)

	id, err := repo.AddOneOffEvent("event-crud original", "one-off description", start, 45, "America/Los_Angeles")
	if err != nil {
		t.Fatalf("AddOneOffEvent() unexpected error = %v", err)
	}
	if id == "" {
		t.Fatalf("AddOneOffEvent() returned an empty id")
	}

	event, err := repo.GetEvent(id)
	if err != nil {
		t.Fatalf("GetEvent() unexpected error = %v", err)
	}
	if event == nil {
		t.Fatalf("GetEvent() returned nil for a just-added event")
	}
	if event.ID != id {
		t.Errorf("event ID = %v, want %v", event.ID, id)
	}
	if event.SeriesID != nil {
		t.Errorf("one-off event SeriesID = %v, want nil", *event.SeriesID)
	}
	if event.OriginalStart != nil {
		t.Errorf("one-off event OriginalStart = %v, want nil", *event.OriginalStart)
	}
	if event.Name != "event-crud original" {
		t.Errorf("event Name = %v, want %v", event.Name, "event-crud original")
	}
	if event.Description != "one-off description" {
		t.Errorf("event Description = %v, want %v", event.Description, "one-off description")
	}
	if !event.StartTime.Equal(start) {
		t.Errorf("event StartTime = %v, want the same instant as %v", event.StartTime, start)
	}
	if event.DurationMinutes != 45 {
		t.Errorf("event DurationMinutes = %v, want 45", event.DurationMinutes)
	}
	if event.Timezone != "America/Los_Angeles" {
		t.Errorf("event Timezone = %v, want America/Los_Angeles", event.Timezone)
	}
	if event.Status != models.ScheduledEventStatusScheduled {
		t.Errorf("new event Status = %v, want %v", event.Status, models.ScheduledEventStatusScheduled)
	}
	if event.FederatedAt != nil {
		t.Errorf("new event FederatedAt = %v, want nil", *event.FederatedAt)
	}
	if event.ReminderSentAt != nil {
		t.Errorf("new event ReminderSentAt = %v, want nil", *event.ReminderSentAt)
	}

	if err := repo.UpdateEventDetails(id, "event-crud renamed", "updated description", 75); err != nil {
		t.Fatalf("UpdateEventDetails() unexpected error = %v", err)
	}
	event, err = repo.GetEvent(id)
	if err != nil || event == nil {
		t.Fatalf("GetEvent() after update = %v, %v", event, err)
	}
	if event.Name != "event-crud renamed" {
		t.Errorf("updated event Name = %v, want %v", event.Name, "event-crud renamed")
	}
	if event.Description != "updated description" {
		t.Errorf("updated event Description = %v, want %v", event.Description, "updated description")
	}
	if event.DurationMinutes != 75 {
		t.Errorf("updated event DurationMinutes = %v, want 75", event.DurationMinutes)
	}
	if !event.StartTime.Equal(start) {
		t.Errorf("UpdateEventDetails() changed StartTime to %v, want %v", event.StartTime, start)
	}
	if event.Timezone != "America/Los_Angeles" {
		t.Errorf("UpdateEventDetails() changed Timezone to %v", event.Timezone)
	}

	if err := repo.DeleteEvent(id); err != nil {
		t.Fatalf("DeleteEvent() unexpected error = %v", err)
	}
	event, err = repo.GetEvent(id)
	if err != nil {
		t.Fatalf("GetEvent() after delete unexpected error = %v", err)
	}
	if event != nil {
		t.Errorf("GetEvent() after delete = %+v, want nil", event)
	}
}

func TestGetEventsInRangeHalfOpen(t *testing.T) {
	repo := Get()

	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("LoadLocation() unexpected error = %v", err)
	}

	from := time.Date(2031, 3, 10, 0, 0, 0, 0, time.UTC)
	to := from.Add(2 * time.Hour)

	atFrom := mustAddOneOffEvent(t, "range at lower bound", from)
	// Same instant as from+1h, expressed in a non-UTC zone: the UTC-bounded
	// query must still find it.
	inside := mustAddOneOffEvent(t, "range inside", from.Add(time.Hour).In(la))
	atTo := mustAddOneOffEvent(t, "range at upper bound", to)
	before := mustAddOneOffEvent(t, "range before lower bound", from.Add(-time.Second))

	events, err := repo.GetEventsInRange(from, to)
	if err != nil {
		t.Fatalf("GetEventsInRange() unexpected error = %v", err)
	}

	if !containsEventID(events, atFrom) {
		t.Errorf("GetEventsInRange() omits the event starting exactly at from (lower bound must be inclusive)")
	}
	if !containsEventID(events, inside) {
		t.Errorf("GetEventsInRange() omits an event inside the range stored from a non-UTC input")
	}
	if containsEventID(events, atTo) {
		t.Errorf("GetEventsInRange() includes the event starting exactly at to (upper bound must be exclusive)")
	}
	if containsEventID(events, before) {
		t.Errorf("GetEventsInRange() includes an event starting before from")
	}
}

func TestGetEventsToFederateSelectsOnce(t *testing.T) {
	repo := Get()

	startingAfter := time.Date(2032, 1, 1, 0, 0, 0, 0, time.UTC)

	pending := mustAddOneOffEvent(t, "federate pending", startingAfter.Add(4*24*time.Hour))
	atBoundary := mustAddOneOffEvent(t, "federate at boundary", startingAfter)
	cancelled := mustAddOneOffEvent(t, "federate cancelled", startingAfter.Add(5*24*time.Hour))
	if err := repo.CancelEvent(cancelled); err != nil {
		t.Fatalf("CancelEvent() unexpected error = %v", err)
	}

	events, err := repo.GetEventsToFederate(startingAfter)
	if err != nil {
		t.Fatalf("GetEventsToFederate() unexpected error = %v", err)
	}
	if !containsEventID(events, pending) {
		t.Errorf("GetEventsToFederate() omits an unfederated scheduled future event")
	}
	if containsEventID(events, atBoundary) {
		t.Errorf("GetEventsToFederate() includes an event starting exactly at the boundary (must be strictly after)")
	}
	if containsEventID(events, cancelled) {
		t.Errorf("GetEventsToFederate() includes a cancelled event")
	}

	stamp := time.Date(2032, 1, 2, 9, 0, 0, 0, time.UTC)
	if err := repo.SetEventFederatedAt(pending, stamp); err != nil {
		t.Fatalf("SetEventFederatedAt() unexpected error = %v", err)
	}
	event, err := repo.GetEvent(pending)
	if err != nil || event == nil {
		t.Fatalf("GetEvent() after federation stamp = %v, %v", event, err)
	}
	if event.FederatedAt == nil {
		t.Fatalf("event FederatedAt = nil after SetEventFederatedAt()")
	}
	if !event.FederatedAt.Equal(stamp) {
		t.Errorf("event FederatedAt = %v, want %v", event.FederatedAt, stamp)
	}

	events, err = repo.GetEventsToFederate(startingAfter)
	if err != nil {
		t.Fatalf("GetEventsToFederate() second call unexpected error = %v", err)
	}
	if containsEventID(events, pending) {
		t.Errorf("GetEventsToFederate() returns an already-federated event a second time")
	}
}

func TestGetEventsNeedingReminderSelectsOnce(t *testing.T) {
	repo := Get()

	startAfter := time.Date(2033, 6, 1, 12, 0, 0, 0, time.UTC)
	startBefore := startAfter.Add(30 * time.Minute)

	atLower := mustAddOneOffEvent(t, "reminder at lower bound", startAfter)
	inside := mustAddOneOffEvent(t, "reminder inside window", startAfter.Add(15*time.Minute))
	atUpper := mustAddOneOffEvent(t, "reminder at upper bound", startBefore)
	pastUpper := mustAddOneOffEvent(t, "reminder past upper bound", startBefore.Add(time.Minute))
	cancelled := mustAddOneOffEvent(t, "reminder cancelled", startAfter.Add(20*time.Minute))
	if err := repo.CancelEvent(cancelled); err != nil {
		t.Fatalf("CancelEvent() unexpected error = %v", err)
	}

	events, err := repo.GetEventsNeedingReminder(startAfter, startBefore)
	if err != nil {
		t.Fatalf("GetEventsNeedingReminder() unexpected error = %v", err)
	}
	if containsEventID(events, atLower) {
		t.Errorf("GetEventsNeedingReminder() includes an event starting exactly at startAfter (lower bound must be exclusive)")
	}
	if !containsEventID(events, inside) {
		t.Errorf("GetEventsNeedingReminder() omits an event inside the window")
	}
	if !containsEventID(events, atUpper) {
		t.Errorf("GetEventsNeedingReminder() omits the event starting exactly at startBefore (upper bound must be inclusive)")
	}
	if containsEventID(events, pastUpper) {
		t.Errorf("GetEventsNeedingReminder() includes an event starting after startBefore")
	}
	if containsEventID(events, cancelled) {
		t.Errorf("GetEventsNeedingReminder() includes a cancelled event")
	}

	stamp := time.Date(2033, 6, 1, 12, 5, 0, 0, time.UTC)
	if err := repo.SetEventReminderSentAt(inside, stamp); err != nil {
		t.Fatalf("SetEventReminderSentAt() unexpected error = %v", err)
	}
	event, err := repo.GetEvent(inside)
	if err != nil || event == nil {
		t.Fatalf("GetEvent() after reminder stamp = %v, %v", event, err)
	}
	if event.ReminderSentAt == nil {
		t.Fatalf("event ReminderSentAt = nil after SetEventReminderSentAt()")
	}
	if !event.ReminderSentAt.Equal(stamp) {
		t.Errorf("event ReminderSentAt = %v, want %v", event.ReminderSentAt, stamp)
	}

	events, err = repo.GetEventsNeedingReminder(startAfter, startBefore)
	if err != nil {
		t.Fatalf("GetEventsNeedingReminder() second call unexpected error = %v", err)
	}
	if containsEventID(events, inside) {
		t.Errorf("GetEventsNeedingReminder() returns an already-reminded event a second time")
	}
	if !containsEventID(events, atUpper) {
		t.Errorf("GetEventsNeedingReminder() dropped an unreminded event after stamping a different one")
	}
}
