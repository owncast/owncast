package handlers

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
	"github.com/sherif-fanous/xmltv"
)

// maxScheduleRangeDays bounds how much schedule a single public request may
// ask for. The viewer calendar asks month-by-month.
const maxScheduleRangeDays = 366

// GetSchedule returns the public schedule of stream events, past and
// upcoming. Always an empty list while the feature is disabled: the schedule
// surface stays inert without leaking whether data exists.
func (h *Handlers) GetSchedule(w http.ResponseWriter, r *http.Request, params generated.GetScheduleParams) {
	if !h.configRepository.GetScheduleEnabled() {
		webutils.WriteResponse(w, []models.ScheduledEvent{})
		return
	}

	now := time.Now()
	from := now.AddDate(0, 0, -7)
	to := now.AddDate(0, 0, 90)
	if params.From != nil {
		from = *params.From
	}
	if params.To != nil {
		to = *params.To
	}

	if !to.After(from) {
		webutils.BadRequestHandler(w, errors.New("to must be after from"))
		return
	}
	if to.Sub(from) > maxScheduleRangeDays*24*time.Hour {
		webutils.BadRequestHandler(w, errors.New("requested schedule range is too large"))
		return
	}

	events, err := h.scheduleEventsRepository.GetEventsInRange(from, to)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	webutils.WriteResponse(w, events)
}

// GetScheduleICS returns a subscribable, always-current iCalendar feed of
// scheduled stream events.
func (h *Handlers) GetScheduleICS(w http.ResponseWriter, r *http.Request) {
	events, err := h.getScheduleFeedEvents()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	middleware.DisableCache(w)
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	_, _ = w.Write([]byte(scheduleICS(events, time.Now())))
}

// GetScheduleXMLTV returns an always-current XMLTV guide of scheduled streams.
func (h *Handlers) GetScheduleXMLTV(w http.ResponseWriter, r *http.Request) {
	events, err := h.getScheduleFeedEvents()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	guide, err := scheduleXMLTV(
		events,
		h.configRepository.GetServerName(),
		h.configRepository.GetServerURL(),
	)
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	middleware.DisableCache(w)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write(guide)
}

func (h *Handlers) getScheduleFeedEvents() ([]models.ScheduledEvent, error) {
	if !h.configRepository.GetScheduleEnabled() {
		return []models.ScheduledEvent{}, nil
	}

	now := time.Now()
	return h.scheduleEventsRepository.GetEventsInRange(now.AddDate(0, 0, -7), now.AddDate(1, 0, 0))
}

func scheduleXMLTV(events []models.ScheduledEvent, serverName string, serverURL string) ([]byte, error) {
	serverName = normalizeXMLTVText(serverName)
	if serverName == "" {
		serverName = "Owncast"
	}

	serverURL = strings.TrimRight(serverURL, "/")
	channelID := "owncast"
	if parsedURL, err := url.Parse(serverURL); err == nil && parsedURL.Hostname() != "" {
		channelID = parsedURL.Hostname()
	}

	generatorName := "Owncast"
	generatorURL := "https://owncast.online"
	guide := xmltv.TV{
		GeneratorInfoName: &generatorName,
		GeneratorInfoURL:  &generatorURL,
		Channels: []xmltv.Channel{
			{
				ID:           channelID,
				DisplayNames: []xmltv.DisplayName{{Text: serverName}},
			},
		},
	}
	if serverURL != "" {
		sourceDataURL := serverURL + "/api/schedule.xml"
		guide.SourceInfoName = &serverName
		guide.SourceInfoURL = &serverURL
		guide.SourceDataURL = &sourceDataURL
		guide.Channels[0].Icons = []xmltv.Icon{{Source: serverURL + "/logo/external"}}
		guide.Channels[0].URLs = []xmltv.URL{{Text: serverURL}}
	}

	for _, event := range events {
		if event.Status == models.ScheduledEventStatusCancelled {
			continue
		}

		start := xmltv.Time{Time: event.StartTime.UTC()}
		stop := xmltv.Time{Time: event.StartTime.UTC().Add(time.Duration(event.DurationMinutes) * time.Minute)}
		programme := xmltv.Programme{
			Start:   start,
			Stop:    &stop,
			Channel: channelID,
			Titles:  []xmltv.Title{{Text: normalizeXMLTVText(event.Name)}},
		}
		if description := normalizeXMLTVText(event.Description); description != "" {
			programme.Descriptions = []xmltv.Description{{Text: description}}
		}
		guide.Programmes = append(guide.Programmes, programme)
	}

	body, err := xml.Marshal(guide)
	if err != nil {
		return nil, err
	}

	return append(append([]byte(xml.Header), body...), '\n'), nil
}

func normalizeXMLTVText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func scheduleICS(events []models.ScheduledEvent, generatedAt time.Time) string {
	var calendar strings.Builder
	for _, line := range []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Owncast//Schedule//EN",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"REFRESH-INTERVAL;VALUE=DURATION:PT15M",
	} {
		writeICSLine(&calendar, line)
	}
	for _, event := range events {
		start := event.StartTime.UTC()
		end := start.Add(time.Duration(event.DurationMinutes) * time.Minute)
		writeICSLine(&calendar, "BEGIN:VEVENT")
		writeICSLine(&calendar, "UID:"+escapeICSText(event.ID)+"@owncast")
		writeICSLine(&calendar, "DTSTAMP:"+formatICSDate(generatedAt))
		writeICSLine(&calendar, "DTSTART:"+formatICSDate(start))
		writeICSLine(&calendar, "DTEND:"+formatICSDate(end))
		writeICSLine(&calendar, "SUMMARY:"+escapeICSText(event.Name))
		if event.Description != "" {
			writeICSLine(&calendar, "DESCRIPTION:"+escapeICSText(event.Description))
		}
		if event.Status == models.ScheduledEventStatusCancelled {
			writeICSLine(&calendar, "STATUS:CANCELLED")
		} else {
			writeICSLine(&calendar, "STATUS:CONFIRMED")
		}
		writeICSLine(&calendar, "END:VEVENT")
	}
	writeICSLine(&calendar, "END:VCALENDAR")
	return calendar.String()
}

func writeICSLine(calendar *strings.Builder, line string) {
	limit := 75
	for len(line) > limit {
		split := limit
		for !utf8.RuneStart(line[split]) {
			split--
		}
		calendar.WriteString(line[:split])
		calendar.WriteString("\r\n ")
		line = line[split:]
		limit = 74
	}
	calendar.WriteString(line)
	calendar.WriteString("\r\n")
}

func formatICSDate(value time.Time) string {
	return value.UTC().Format("20060102T150405Z")
}

func escapeICSText(value string) string {
	return strings.NewReplacer(
		`\`, `\\`,
		"\r\n", `\n`,
		"\n", `\n`,
		"\r", `\n`,
		";", `\;`,
		",", `\,`,
	).Replace(value)
}
