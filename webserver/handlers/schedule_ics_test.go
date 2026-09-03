package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
)

func TestScheduleICS(t *testing.T) {
	start := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.FixedZone("EDT", -4*60*60))
	feed := scheduleICS([]models.ScheduledEvent{
		{
			ID:              "event-1",
			Name:            "Show; one",
			Description:     "Line one,\\line two\nline three",
			StartTime:       start,
			DurationMinutes: 90,
			Status:          models.ScheduledEventStatusScheduled,
		},
		{
			ID:              "event-2",
			Name:            "Cancelled show",
			StartTime:       start.Add(24 * time.Hour),
			DurationMinutes: 30,
			Status:          models.ScheduledEventStatusCancelled,
		},
	}, time.Date(2030, time.September, 1, 18, 0, 0, 0, time.UTC))

	for _, expected := range []string{
		"BEGIN:VCALENDAR\r\n",
		"UID:event-1@owncast\r\n",
		"DTSTAMP:20300901T180000Z\r\n",
		"DTSTART:20300908T220000Z\r\n",
		"DTEND:20300908T233000Z\r\n",
		"SUMMARY:Show\\; one\r\n",
		"DESCRIPTION:Line one\\,\\\\line two\\nline three\r\n",
		"STATUS:CANCELLED\r\n",
		"END:VCALENDAR\r\n",
	} {
		if !strings.Contains(feed, expected) {
			t.Errorf("feed missing %q", expected)
		}
	}
}

func TestWriteICSLineFoldsUTF8(t *testing.T) {
	line := "DESCRIPTION:" + strings.Repeat("é", 50)
	var output strings.Builder
	writeICSLine(&output, line)

	for _, physicalLine := range strings.Split(strings.TrimSuffix(output.String(), "\r\n"), "\r\n") {
		if len([]byte(physicalLine)) > 75 {
			t.Errorf("folded line is %d bytes, expected at most 75", len([]byte(physicalLine)))
		}
	}
	if unfolded := strings.ReplaceAll(strings.TrimSuffix(output.String(), "\r\n"), "\r\n ", ""); unfolded != line {
		t.Errorf("folding changed content: %q", unfolded)
	}
}

func TestGetScheduleICSReturnsCurrentSchedule(t *testing.T) {
	configRepo := configrepository.New(testDatastore)
	if err := configRepo.SetScheduleEnabled(true); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = configRepo.SetScheduleEnabled(false)
	})

	eventsRepo := scheduleeventsrepository.New(testDatastore)
	eventID, err := eventsRepo.AddOneOffEvent(
		"Original name",
		"Description",
		time.Now().Add(24*time.Hour),
		60,
		"UTC",
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = eventsRepo.DeleteEvent(eventID)
	})

	h := *testHandlers
	h.configRepository = configRepo
	h.scheduleEventsRepository = eventsRepo
	handler := New(&h).Handler()

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/schedule.ics", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", first.Code)
	}
	if contentType := first.Header().Get("Content-Type"); contentType != "text/calendar; charset=utf-8" {
		t.Errorf("unexpected Content-Type %q", contentType)
	}
	if cacheControl := first.Header().Get("Cache-Control"); !strings.Contains(cacheControl, "no-cache") {
		t.Errorf("expected no-cache response, got %q", cacheControl)
	}
	if !strings.Contains(first.Body.String(), "SUMMARY:Original name\r\n") {
		t.Fatal("initial feed did not contain the scheduled event")
	}

	if err := eventsRepo.UpdateEventDetails(eventID, "Updated name", "Description", 60); err != nil {
		t.Fatal(err)
	}
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/schedule.ics", nil))
	if !strings.Contains(second.Body.String(), "SUMMARY:Updated name\r\n") {
		t.Error("feed did not reflect the updated event")
	}
	if strings.Contains(second.Body.String(), "SUMMARY:Original name\r\n") {
		t.Error("feed returned stale event data")
	}
}
