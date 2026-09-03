package handlers

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	"github.com/sherif-fanous/xmltv"
)

func TestScheduleXMLTV(t *testing.T) {
	start := time.Date(2030, time.September, 8, 18, 0, 0, 0, time.FixedZone("EDT", -4*60*60))
	feed, err := scheduleXMLTV([]models.ScheduledEvent{
		{
			ID:              "event-1",
			Name:            "Show & one",
			Description:     "Line one\nline <two>",
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
	}, "My & Channel\nName", "https://stream.example/")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(feed, []byte(xml.Header)) {
		t.Fatal("guide is missing the XML declaration")
	}

	var guide xmltv.TV
	if err := xml.Unmarshal(feed, &guide); err != nil {
		t.Fatalf("guide is not valid XMLTV: %v", err)
	}
	if len(guide.Channels) != 1 {
		t.Fatalf("got %d channels, want 1", len(guide.Channels))
	}
	channel := guide.Channels[0]
	if channel.ID != "stream.example" {
		t.Errorf("channel ID = %q, want stream.example", channel.ID)
	}
	if got := channel.DisplayNames[0].Text; got != "My & Channel Name" {
		t.Errorf("display name = %q, want normalized server name", got)
	}
	if got := channel.Icons[0].Source; got != "https://stream.example/logo/external" {
		t.Errorf("channel icon = %q", got)
	}
	if guide.SourceDataURL == nil || *guide.SourceDataURL != "https://stream.example/api/schedule.xml" {
		t.Errorf("source data URL = %v", guide.SourceDataURL)
	}

	if len(guide.Programmes) != 1 {
		t.Fatalf("got %d programmes, want the one non-cancelled event", len(guide.Programmes))
	}
	programme := guide.Programmes[0]
	if programme.Channel != channel.ID {
		t.Errorf("programme channel = %q, want %q", programme.Channel, channel.ID)
	}
	if got := programme.Start.UTC(); !got.Equal(time.Date(2030, time.September, 8, 22, 0, 0, 0, time.UTC)) {
		t.Errorf("programme start = %v", got)
	}
	if programme.Stop == nil || !programme.Stop.UTC().Equal(time.Date(2030, time.September, 8, 23, 30, 0, 0, time.UTC)) {
		t.Errorf("programme stop = %v", programme.Stop)
	}
	if got := programme.Titles[0].Text; got != "Show & one" {
		t.Errorf("programme title = %q", got)
	}
	if got := programme.Descriptions[0].Text; got != "Line one line <two>" {
		t.Errorf("programme description = %q", got)
	}
}

func TestGetScheduleXMLTV(t *testing.T) {
	configRepo := configrepository.New(testDatastore)
	if err := configRepo.SetScheduleEnabled(true); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = configRepo.SetScheduleEnabled(false)
	})

	eventsRepo := scheduleeventsrepository.New(testDatastore)
	eventID, err := eventsRepo.AddOneOffEvent(
		"Guide event",
		"Description",
		"",
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

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/schedule.xml", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/xml; charset=utf-8" {
		t.Errorf("unexpected Content-Type %q", contentType)
	}
	if cacheControl := response.Header().Get("Cache-Control"); !strings.Contains(cacheControl, "no-cache") {
		t.Errorf("expected no-cache response, got %q", cacheControl)
	}

	var guide xmltv.TV
	if err := xml.Unmarshal(response.Body.Bytes(), &guide); err != nil {
		t.Fatalf("response is not valid XMLTV: %v", err)
	}
	if len(guide.Programmes) != 1 || guide.Programmes[0].Titles[0].Text != "Guide event" {
		t.Fatalf("unexpected programmes: %+v", guide.Programmes)
	}

	if err := configRepo.SetScheduleEnabled(false); err != nil {
		t.Fatal(err)
	}
	disabledResponse := httptest.NewRecorder()
	handler.ServeHTTP(disabledResponse, httptest.NewRequest(http.MethodGet, "/schedule.xml", nil))
	guide = xmltv.TV{}
	if err := xml.Unmarshal(disabledResponse.Body.Bytes(), &guide); err != nil {
		t.Fatalf("disabled response is not valid XMLTV: %v", err)
	}
	if len(guide.Programmes) != 0 {
		t.Errorf("disabled schedule returned %d programmes", len(guide.Programmes))
	}
}
