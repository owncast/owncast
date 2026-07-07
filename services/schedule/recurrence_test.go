package schedule

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/migrations"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	"github.com/owncast/owncast/services/datastore"
)

var testRepo scheduleeventsrepository.ScheduleEventsRepository

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

	testRepo = scheduleeventsrepository.New(ds)
	scheduleeventsrepository.SetGlobalInstance(testRepo)

	m.Run()
}

// mustAddSeries creates a series through the repository and deactivates it at
// test end so MaterializeAllSeries in a later test only sees its own series
// in the shared datastore.
func mustAddSeries(t *testing.T, name, recurrence string, durationMinutes int) models.ScheduledEventSeries {
	t.Helper()
	id, err := testRepo.AddSeries(name, "description for "+name, recurrence, durationMinutes)
	if err != nil {
		t.Fatalf("AddSeries(%s) unexpected error = %v", name, err)
	}
	t.Cleanup(func() {
		if err := testRepo.SetSeriesActive(id, false); err != nil {
			t.Errorf("cleanup SetSeriesActive(%s) unexpected error = %v", id, err)
		}
	})
	series, err := testRepo.GetSeries(id)
	if err != nil {
		t.Fatalf("GetSeries(%s) unexpected error = %v", id, err)
	}
	if series == nil {
		t.Fatalf("GetSeries(%s) returned nil for a just-added series", id)
	}
	return *series
}

func mustGetEventsForSeries(t *testing.T, seriesID string) []models.ScheduledEvent {
	t.Helper()
	events, err := testRepo.GetEventsForSeries(seriesID)
	if err != nil {
		t.Fatalf("GetEventsForSeries(%s) unexpected error = %v", seriesID, err)
	}
	return events
}

func TestWeeklyRecurrenceRoundTrip(t *testing.T) {
	recurrence := "DTSTART;TZID=America/New_York:20270601T190000\nRRULE:FREQ=WEEKLY;BYDAY=TU"

	set, err := ParseRecurrence(recurrence)
	if err != nil {
		t.Fatalf("ParseRecurrence() unexpected error = %v", err)
	}
	if timezone := RecurrenceTimezone(set); timezone != "America/New_York" {
		t.Errorf("RecurrenceTimezone() = %v, want America/New_York", timezone)
	}

	from := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 6, 15, 0, 0, 0, 0, time.UTC)
	occurrences, err := ExpandBetween(recurrence, from, to)
	if err != nil {
		t.Fatalf("ExpandBetween() unexpected error = %v", err)
	}

	// Tuesdays June 1 and June 8 at 19:00 EDT (UTC-4). June 15's occurrence
	// starts at 23:00Z, past the half-open upper bound of June 15 00:00Z.
	expected := []time.Time{
		time.Date(2027, 6, 1, 23, 0, 0, 0, time.UTC),
		time.Date(2027, 6, 8, 23, 0, 0, 0, time.UTC),
	}
	if len(occurrences) != len(expected) {
		t.Fatalf("ExpandBetween() returned %d occurrences %v, want %d", len(occurrences), occurrences, len(expected))
	}
	for i, want := range expected {
		if !occurrences[i].Equal(want) {
			t.Errorf("occurrence[%d] = %v, want %v", i, occurrences[i], want)
		}
	}
}

func TestDSTSpringForwardExactInstants(t *testing.T) {
	recurrence := "DTSTART;TZID=America/Los_Angeles:20270301T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO"

	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("LoadLocation() unexpected error = %v", err)
	}

	from := time.Date(2027, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 4, 1, 0, 0, 0, 0, time.UTC)
	occurrences, err := ExpandBetween(recurrence, from, to)
	if err != nil {
		t.Fatalf("ExpandBetween() unexpected error = %v", err)
	}

	// US spring-forward is Sunday 2027-03-14. Mondays before it are 18:00
	// PST (UTC-8), Mondays after it are 18:00 PDT (UTC-7): the wall clock
	// holds at 18:00 while the UTC instant shifts by an hour.
	expected := []time.Time{
		time.Date(2027, 3, 2, 2, 0, 0, 0, time.UTC),  // Mon Mar 1 18:00 PST
		time.Date(2027, 3, 9, 2, 0, 0, 0, time.UTC),  // Mon Mar 8 18:00 PST
		time.Date(2027, 3, 16, 1, 0, 0, 0, time.UTC), // Mon Mar 15 18:00 PDT
		time.Date(2027, 3, 23, 1, 0, 0, 0, time.UTC), // Mon Mar 22 18:00 PDT
		time.Date(2027, 3, 30, 1, 0, 0, 0, time.UTC), // Mon Mar 29 18:00 PDT
	}
	if len(occurrences) != len(expected) {
		t.Fatalf("ExpandBetween() returned %d occurrences %v, want %d", len(occurrences), occurrences, len(expected))
	}
	for i, want := range expected {
		if !occurrences[i].Equal(want) {
			t.Errorf("occurrence[%d] = %v, want the exact instant %v", i, occurrences[i], want)
		}
		local := occurrences[i].In(la)
		if local.Hour() != 18 {
			t.Errorf("occurrence[%d] wall-clock hour in America/Los_Angeles = %d, want 18", i, local.Hour())
		}
		if local.Weekday() != time.Monday {
			t.Errorf("occurrence[%d] weekday in America/Los_Angeles = %v, want Monday", i, local.Weekday())
		}
	}
}

func TestParseRecurrenceRejectsInvalidValues(t *testing.T) {
	cases := []struct {
		name       string
		recurrence string
	}{
		{"empty", ""},
		{"garbage", "this is not a recurrence rule"},
		{"missing rrule", "DTSTART;TZID=America/Los_Angeles:20270301T180000"},
		{"missing dtstart", "RRULE:FREQ=WEEKLY;BYDAY=MO"},
		{"unknown timezone", "DTSTART;TZID=Mars/Olympus:20270301T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			set, err := ParseRecurrence(tc.recurrence)
			if err == nil {
				t.Errorf("ParseRecurrence(%q) = nil error, want an error", tc.recurrence)
			}
			if set != nil {
				t.Errorf("ParseRecurrence(%q) returned a non-nil set alongside the error", tc.recurrence)
			}
		})
	}
}

func TestExpandBetweenGuardsAgainstDenseRules(t *testing.T) {
	t.Run("occurrence cap inside the window", func(t *testing.T) {
		// A minutely rule over 30 days is 43200 occurrences, far past the
		// per-expansion cap.
		recurrence := "DTSTART;TZID=UTC:20260101T000000\nRRULE:FREQ=MINUTELY"
		from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		to := from.Add(30 * 24 * time.Hour)
		occurrences, err := ExpandBetween(recurrence, from, to)
		if err == nil {
			t.Fatalf("ExpandBetween() = nil error, want the occurrence-cap error")
		}
		if !strings.Contains(err.Error(), "too many occurrences") {
			t.Errorf("ExpandBetween() error = %v, want it to mention too many occurrences", err)
		}
		if occurrences != nil {
			t.Errorf("ExpandBetween() returned %d occurrences alongside the error, want nil", len(occurrences))
		}
	})

	t.Run("iteration cap before a distant window", func(t *testing.T) {
		// A secondly rule whose DTSTART is decades before the window would
		// walk over a billion instants before producing a single result. The
		// iteration cap must refuse it promptly instead.
		recurrence := "DTSTART;TZID=UTC:20000101T000000\nRRULE:FREQ=SECONDLY"
		from := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
		to := from.Add(24 * time.Hour)
		started := time.Now()
		occurrences, err := ExpandBetween(recurrence, from, to)
		elapsed := time.Since(started)
		if err == nil {
			t.Fatalf("ExpandBetween() = nil error, want the iteration-cap error")
		}
		if !strings.Contains(err.Error(), "too dense") {
			t.Errorf("ExpandBetween() error = %v, want it to mention the rule being too dense", err)
		}
		if occurrences != nil {
			t.Errorf("ExpandBetween() returned %d occurrences alongside the error, want nil", len(occurrences))
		}
		if elapsed > 10*time.Second {
			t.Errorf("ExpandBetween() took %v to hit the iteration cap, want a prompt refusal", elapsed)
		}
	})

	t.Run("ordinary weekly rule is untouched by the guards", func(t *testing.T) {
		// Mondays Jan 5, 12, 19 2026 inside [Jan 1, Jan 22).
		recurrence := "DTSTART;TZID=UTC:20260105T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO"
		from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		to := from.Add(21 * 24 * time.Hour)
		occurrences, err := ExpandBetween(recurrence, from, to)
		if err != nil {
			t.Fatalf("ExpandBetween() unexpected error = %v", err)
		}
		if len(occurrences) != 3 {
			t.Errorf("ExpandBetween() returned %d occurrences %v, want 3", len(occurrences), occurrences)
		}
	})
}

func TestMaterializeSeriesIsIdempotent(t *testing.T) {
	// Mondays Jan 3, 10, 17 2028 at 12:00 UTC inside [Jan 1, Jan 22).
	series := mustAddSeries(t, "materialize idempotent", "DTSTART;TZID=UTC:20280103T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)
	now := time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC)
	horizon := 21 * 24 * time.Hour

	inserted, err := MaterializeSeries(testRepo, series, now, horizon)
	if err != nil {
		t.Fatalf("MaterializeSeries() unexpected error = %v", err)
	}
	if inserted != 3 {
		t.Fatalf("MaterializeSeries() inserted = %d, want 3", inserted)
	}

	events := mustGetEventsForSeries(t, series.ID)
	if len(events) != 3 {
		t.Fatalf("GetEventsForSeries() returned %d events, want 3", len(events))
	}
	expectedStarts := []time.Time{
		time.Date(2028, 1, 3, 12, 0, 0, 0, time.UTC),
		time.Date(2028, 1, 10, 12, 0, 0, 0, time.UTC),
		time.Date(2028, 1, 17, 12, 0, 0, 0, time.UTC),
	}
	for i, want := range expectedStarts {
		if !events[i].StartTime.Equal(want) {
			t.Errorf("event[%d] StartTime = %v, want %v", i, events[i].StartTime, want)
		}
		if events[i].OriginalStart == nil || !events[i].OriginalStart.Equal(want) {
			t.Errorf("event[%d] OriginalStart = %v, want %v", i, events[i].OriginalStart, want)
		}
		if events[i].SeriesID == nil || *events[i].SeriesID != series.ID {
			t.Errorf("event[%d] SeriesID = %v, want %v", i, events[i].SeriesID, series.ID)
		}
		if events[i].Name != series.Name {
			t.Errorf("event[%d] Name = %v, want %v", i, events[i].Name, series.Name)
		}
		if events[i].Timezone != "UTC" {
			t.Errorf("event[%d] Timezone = %v, want UTC", i, events[i].Timezone)
		}
	}

	insertedAgain, err := MaterializeSeries(testRepo, series, now, horizon)
	if err != nil {
		t.Fatalf("second MaterializeSeries() unexpected error = %v", err)
	}
	if insertedAgain != 0 {
		t.Errorf("second MaterializeSeries() inserted = %d, want 0", insertedAgain)
	}

	eventsAfter := mustGetEventsForSeries(t, series.ID)
	if len(eventsAfter) != len(events) {
		t.Fatalf("row count changed after re-materialization: %d, want %d", len(eventsAfter), len(events))
	}
	for i := range events {
		if eventsAfter[i].ID != events[i].ID {
			t.Errorf("event[%d] ID changed after re-materialization: %v, want %v", i, eventsAfter[i].ID, events[i].ID)
		}
	}
}

func TestCancelledOccurrenceIsNotResurrected(t *testing.T) {
	// Mondays Feb 7, 14, 21 2028 at 12:00 UTC inside [Feb 1, Feb 22).
	series := mustAddSeries(t, "cancel lifecycle", "DTSTART;TZID=UTC:20280207T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)
	now := time.Date(2028, 2, 1, 0, 0, 0, 0, time.UTC)
	horizon := 21 * 24 * time.Hour

	if inserted, err := MaterializeSeries(testRepo, series, now, horizon); err != nil || inserted != 3 {
		t.Fatalf("MaterializeSeries() = %d, %v, want 3, nil", inserted, err)
	}
	events := mustGetEventsForSeries(t, series.ID)
	if len(events) != 3 {
		t.Fatalf("GetEventsForSeries() returned %d events, want 3", len(events))
	}
	target := events[1]

	if err := testRepo.CancelEvent(target.ID); err != nil {
		t.Fatalf("CancelEvent() unexpected error = %v", err)
	}
	cancelled, err := testRepo.GetEvent(target.ID)
	if err != nil {
		t.Fatalf("GetEvent() unexpected error = %v", err)
	}
	if cancelled == nil {
		t.Fatalf("CancelEvent() deleted the row, want it kept")
	}
	if cancelled.Status != models.ScheduledEventStatusCancelled {
		t.Errorf("cancelled event Status = %v, want %v", cancelled.Status, models.ScheduledEventStatusCancelled)
	}
	if !cancelled.StartTime.Equal(target.StartTime) {
		t.Errorf("CancelEvent() changed StartTime to %v, want %v", cancelled.StartTime, target.StartTime)
	}

	inserted, err := MaterializeSeries(testRepo, series, now, horizon)
	if err != nil {
		t.Fatalf("MaterializeSeries() after cancel unexpected error = %v", err)
	}
	if inserted != 0 {
		t.Errorf("MaterializeSeries() after cancel inserted = %d, want 0", inserted)
	}

	eventsAfter := mustGetEventsForSeries(t, series.ID)
	if len(eventsAfter) != 3 {
		t.Fatalf("row count after re-materialization = %d, want 3", len(eventsAfter))
	}
	slotCount := 0
	for _, e := range eventsAfter {
		if e.OriginalStart != nil && e.OriginalStart.Equal(*target.OriginalStart) {
			slotCount++
			if e.ID != target.ID {
				t.Errorf("cancelled slot re-created with a new row %v, want only %v", e.ID, target.ID)
			}
			if e.Status != models.ScheduledEventStatusCancelled {
				t.Errorf("cancelled slot Status = %v after re-materialization, want %v", e.Status, models.ScheduledEventStatusCancelled)
			}
		}
	}
	if slotCount != 1 {
		t.Errorf("cancelled slot occupies %d rows after re-materialization, want 1", slotCount)
	}
}

func TestMovedOccurrenceKeepsSlotIdentity(t *testing.T) {
	// Mondays Mar 6, 13, 20 2028 at 12:00 UTC inside [Mar 1, Mar 22).
	series := mustAddSeries(t, "move lifecycle", "DTSTART;TZID=UTC:20280306T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)
	now := time.Date(2028, 3, 1, 0, 0, 0, 0, time.UTC)
	horizon := 21 * 24 * time.Hour

	if inserted, err := MaterializeSeries(testRepo, series, now, horizon); err != nil || inserted != 3 {
		t.Fatalf("MaterializeSeries() = %d, %v, want 3, nil", inserted, err)
	}
	events := mustGetEventsForSeries(t, series.ID)
	if len(events) != 3 {
		t.Fatalf("GetEventsForSeries() returned %d events, want 3", len(events))
	}
	target := events[1]
	originalStart := *target.OriginalStart
	newStart := target.StartTime.Add(2 * time.Hour)

	if err := testRepo.MoveEvent(target.ID, newStart); err != nil {
		t.Fatalf("MoveEvent() unexpected error = %v", err)
	}
	moved, err := testRepo.GetEvent(target.ID)
	if err != nil || moved == nil {
		t.Fatalf("GetEvent() after move = %v, %v", moved, err)
	}
	if !moved.StartTime.Equal(newStart) {
		t.Errorf("moved event StartTime = %v, want %v", moved.StartTime, newStart)
	}
	if moved.OriginalStart == nil || !moved.OriginalStart.Equal(originalStart) {
		t.Errorf("MoveEvent() changed OriginalStart to %v, want %v", moved.OriginalStart, originalStart)
	}

	inserted, err := MaterializeSeries(testRepo, series, now, horizon)
	if err != nil {
		t.Fatalf("MaterializeSeries() after move unexpected error = %v", err)
	}
	if inserted != 0 {
		t.Errorf("MaterializeSeries() after move inserted = %d, want 0 (vacated slot must not be re-inserted)", inserted)
	}

	eventsAfter := mustGetEventsForSeries(t, series.ID)
	if len(eventsAfter) != 3 {
		t.Fatalf("row count after re-materialization = %d, want 3", len(eventsAfter))
	}
	slotCount := 0
	for _, e := range eventsAfter {
		if e.OriginalStart != nil && e.OriginalStart.Equal(originalStart) {
			slotCount++
			if !e.StartTime.Equal(newStart) {
				t.Errorf("moved slot StartTime = %v after re-materialization, want %v", e.StartTime, newStart)
			}
		}
	}
	if slotCount != 1 {
		t.Errorf("moved slot occupies %d rows after re-materialization, want 1", slotCount)
	}
}

func TestRegenerateSeriesLeavesFederatedRowsAlone(t *testing.T) {
	// Mondays Apr 3, 10, 17 2028 at 12:00 UTC inside [Apr 1, Apr 22).
	recurrence := "DTSTART;TZID=UTC:20280403T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO"
	series := mustAddSeries(t, "regenerate original", recurrence, 60)
	now := time.Date(2028, 4, 1, 0, 0, 0, 0, time.UTC)
	horizon := 21 * 24 * time.Hour

	if inserted, err := MaterializeSeries(testRepo, series, now, horizon); err != nil || inserted != 3 {
		t.Fatalf("MaterializeSeries() = %d, %v, want 3, nil", inserted, err)
	}
	events := mustGetEventsForSeries(t, series.ID)
	if len(events) != 3 {
		t.Fatalf("GetEventsForSeries() returned %d events, want 3", len(events))
	}
	federated := events[1]
	unfederatedIDs := map[string]bool{events[0].ID: true, events[2].ID: true}

	stamp := time.Date(2028, 4, 1, 6, 0, 0, 0, time.UTC)
	if err := testRepo.SetEventFederatedAt(federated.ID, stamp); err != nil {
		t.Fatalf("SetEventFederatedAt() unexpected error = %v", err)
	}

	if err := testRepo.UpdateSeries(series.ID, "regenerate renamed", "updated description", recurrence, 90); err != nil {
		t.Fatalf("UpdateSeries() unexpected error = %v", err)
	}
	updated, err := testRepo.GetSeries(series.ID)
	if err != nil || updated == nil {
		t.Fatalf("GetSeries() after update = %v, %v", updated, err)
	}
	if updated.Name != "regenerate renamed" {
		t.Fatalf("updated series Name = %v, want %v", updated.Name, "regenerate renamed")
	}

	inserted, err := RegenerateSeries(testRepo, *updated, now, horizon)
	if err != nil {
		t.Fatalf("RegenerateSeries() unexpected error = %v", err)
	}
	if inserted != 2 {
		t.Errorf("RegenerateSeries() inserted = %d, want 2 (only the unfederated slots)", inserted)
	}

	after := mustGetEventsForSeries(t, series.ID)
	if len(after) != 3 {
		t.Fatalf("row count after regeneration = %d, want 3", len(after))
	}
	for _, e := range after {
		if e.OriginalStart == nil {
			t.Fatalf("regenerated event %v has nil OriginalStart", e.ID)
		}
		if e.OriginalStart.Equal(*federated.OriginalStart) {
			// The announced row survives untouched.
			if e.ID != federated.ID {
				t.Errorf("federated row replaced: ID = %v, want %v", e.ID, federated.ID)
			}
			if e.Name != series.Name {
				t.Errorf("federated row Name = %v, want the pre-edit name %v", e.Name, series.Name)
			}
			if e.FederatedAt == nil || !e.FederatedAt.Equal(stamp) {
				t.Errorf("federated row FederatedAt = %v, want %v", e.FederatedAt, stamp)
			}
		} else {
			// Unannounced rows are fresh rows from the edited series.
			if unfederatedIDs[e.ID] {
				t.Errorf("unfederated row %v kept its old ID through regeneration, want a fresh row", e.ID)
			}
			if e.Name != "regenerate renamed" {
				t.Errorf("regenerated row Name = %v, want %v", e.Name, "regenerate renamed")
			}
			if e.DurationMinutes != 90 {
				t.Errorf("regenerated row DurationMinutes = %d, want 90", e.DurationMinutes)
			}
		}
	}
}

func TestRegenerateSeriesLeavesCancelledRowsAlone(t *testing.T) {
	// Mondays May 1, 8, 15 2028 at 12:00 UTC inside [Apr 26, May 17). Every
	// row stays unfederated so only the cancelled status protects the slot.
	recurrence := "DTSTART;TZID=UTC:20280501T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO"
	series := mustAddSeries(t, "regenerate cancelled original", recurrence, 60)
	now := time.Date(2028, 4, 26, 0, 0, 0, 0, time.UTC)
	horizon := 21 * 24 * time.Hour

	if inserted, err := MaterializeSeries(testRepo, series, now, horizon); err != nil || inserted != 3 {
		t.Fatalf("MaterializeSeries() = %d, %v, want 3, nil", inserted, err)
	}
	events := mustGetEventsForSeries(t, series.ID)
	if len(events) != 3 {
		t.Fatalf("GetEventsForSeries() returned %d events, want 3", len(events))
	}
	cancelled := events[1]
	if err := testRepo.CancelEvent(cancelled.ID); err != nil {
		t.Fatalf("CancelEvent() unexpected error = %v", err)
	}

	if err := testRepo.UpdateSeries(series.ID, "regenerate cancelled renamed", "updated description", recurrence, 90); err != nil {
		t.Fatalf("UpdateSeries() unexpected error = %v", err)
	}
	updated, err := testRepo.GetSeries(series.ID)
	if err != nil || updated == nil {
		t.Fatalf("GetSeries() after update = %v, %v", updated, err)
	}

	inserted, err := RegenerateSeries(testRepo, *updated, now, horizon)
	if err != nil {
		t.Fatalf("RegenerateSeries() unexpected error = %v", err)
	}
	if inserted != 2 {
		t.Errorf("RegenerateSeries() inserted = %d, want 2 (the cancelled slot must stay occupied)", inserted)
	}

	after := mustGetEventsForSeries(t, series.ID)
	if len(after) != 3 {
		t.Fatalf("row count after regeneration = %d, want 3", len(after))
	}
	cancelledSlotRows := 0
	for _, e := range after {
		if e.OriginalStart == nil {
			t.Fatalf("event %v has nil OriginalStart after regeneration", e.ID)
		}
		if e.OriginalStart.Equal(*cancelled.OriginalStart) {
			cancelledSlotRows++
			// The cancelled row survives as the same row, untouched by the
			// series edit, and its slot is not re-filled with a fresh
			// scheduled occurrence.
			if e.ID != cancelled.ID {
				t.Errorf("cancelled row replaced: ID = %v, want %v", e.ID, cancelled.ID)
			}
			if e.Status != models.ScheduledEventStatusCancelled {
				t.Errorf("cancelled row Status = %v after regeneration, want %v", e.Status, models.ScheduledEventStatusCancelled)
			}
			if e.Name != series.Name {
				t.Errorf("cancelled row Name = %v, want the pre-edit name %v", e.Name, series.Name)
			}
		} else {
			if e.Name != "regenerate cancelled renamed" {
				t.Errorf("regenerated row Name = %v, want %v", e.Name, "regenerate cancelled renamed")
			}
		}
	}
	if cancelledSlotRows != 1 {
		t.Errorf("cancelled slot occupies %d rows after regeneration, want exactly 1", cancelledSlotRows)
	}
}

func TestMaterializeHorizonIsHalfOpen(t *testing.T) {
	// Monday 2029-05-07. now+horizon lands exactly on the first series'
	// occurrence instant, one minute past the second's.
	occurrenceAtBound := time.Date(2029, 5, 7, 18, 0, 0, 0, time.UTC)
	now := occurrenceAtBound.Add(-8 * time.Hour)
	horizon := 8 * time.Hour

	atBound := mustAddSeries(t, "horizon at bound", "DTSTART;TZID=UTC:20290507T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)
	insideBound := mustAddSeries(t, "horizon inside bound", "DTSTART;TZID=UTC:20290507T175900\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)

	inserted, err := MaterializeSeries(testRepo, atBound, now, horizon)
	if err != nil {
		t.Fatalf("MaterializeSeries(atBound) unexpected error = %v", err)
	}
	if inserted != 0 {
		t.Errorf("occurrence exactly at now+horizon was materialized (%d inserts), want 0", inserted)
	}
	if events := mustGetEventsForSeries(t, atBound.ID); len(events) != 0 {
		t.Errorf("series at bound has %d rows, want 0", len(events))
	}

	inserted, err = MaterializeSeries(testRepo, insideBound, now, horizon)
	if err != nil {
		t.Fatalf("MaterializeSeries(insideBound) unexpected error = %v", err)
	}
	if inserted != 1 {
		t.Fatalf("occurrence at horizon-1min inserted = %d, want 1", inserted)
	}
	events := mustGetEventsForSeries(t, insideBound.ID)
	if len(events) != 1 {
		t.Fatalf("series inside bound has %d rows, want 1", len(events))
	}
	if want := occurrenceAtBound.Add(-time.Minute); !events[0].StartTime.Equal(want) {
		t.Errorf("materialized StartTime = %v, want %v", events[0].StartTime, want)
	}
}

func TestMaterializeSeriesZeroHorizon(t *testing.T) {
	// Monday 2029-06-04 at 18:00 UTC.
	series := mustAddSeries(t, "zero horizon", "DTSTART;TZID=UTC:20290604T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)
	now := time.Date(2029, 6, 1, 0, 0, 0, 0, time.UTC)

	inserted, err := MaterializeSeries(testRepo, series, now, 0)
	if err != nil {
		t.Fatalf("MaterializeSeries() with zero horizon unexpected error = %v", err)
	}
	if inserted != 0 {
		t.Errorf("MaterializeSeries() with zero horizon inserted = %d, want 0", inserted)
	}
	if events := mustGetEventsForSeries(t, series.ID); len(events) != 0 {
		t.Errorf("zero-horizon series has %d rows, want 0", len(events))
	}
}

func TestMaterializeAllSeriesSkipsCorruptSeries(t *testing.T) {
	// Mon May 1 and Tue May 2 2028 at 12:00 UTC inside [May 1, May 8).
	valid1 := mustAddSeries(t, "materialize-all valid monday", "DTSTART;TZID=UTC:20280501T120000\nRRULE:FREQ=WEEKLY;BYDAY=MO", 60)
	corrupt := mustAddSeries(t, "materialize-all corrupt", "this is not a recurrence rule", 60)
	valid2 := mustAddSeries(t, "materialize-all valid tuesday", "DTSTART;TZID=UTC:20280502T120000\nRRULE:FREQ=WEEKLY;BYDAY=TU", 60)

	now := time.Date(2028, 5, 1, 0, 0, 0, 0, time.UTC)
	horizon := 7 * 24 * time.Hour

	inserted, err := MaterializeAllSeries(testRepo, now, horizon)
	if err != nil {
		t.Fatalf("MaterializeAllSeries() unexpected error = %v", err)
	}
	if inserted != 2 {
		t.Errorf("MaterializeAllSeries() inserted = %d, want 2 (one per valid series)", inserted)
	}

	events1 := mustGetEventsForSeries(t, valid1.ID)
	if len(events1) != 1 {
		t.Fatalf("valid series 1 has %d rows, want 1", len(events1))
	}
	if want := time.Date(2028, 5, 1, 12, 0, 0, 0, time.UTC); !events1[0].StartTime.Equal(want) {
		t.Errorf("valid series 1 StartTime = %v, want %v", events1[0].StartTime, want)
	}

	events2 := mustGetEventsForSeries(t, valid2.ID)
	if len(events2) != 1 {
		t.Fatalf("valid series 2 has %d rows, want 1", len(events2))
	}
	if want := time.Date(2028, 5, 2, 12, 0, 0, 0, time.UTC); !events2[0].StartTime.Equal(want) {
		t.Errorf("valid series 2 StartTime = %v, want %v", events2[0].StartTime, want)
	}

	if events := mustGetEventsForSeries(t, corrupt.ID); len(events) != 0 {
		t.Errorf("corrupt series has %d rows, want 0", len(events))
	}
}
