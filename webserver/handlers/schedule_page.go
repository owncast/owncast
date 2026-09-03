package handlers

import (
	"bytes"
	"fmt"
	"html"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
)

// SchedulePageHandler serves the statically-exported schedule page for every
// event URL and adds event metadata for link previews.
func (h *Handlers) SchedulePageHandler(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/schedule/"), "/")
	if decodedID, err := url.PathUnescape(eventID); err == nil {
		eventID = decodedID
	}

	if eventID != "" && h.scheduleEventsRepository != nil {
		if event, err := h.scheduleEventsRepository.GetEvent(eventID); err == nil && event != nil {
			if h.serveScheduleEventPage(w, r, event) {
				return
			}
		}
	}

	r.URL.Path = "/schedule/"
	serveWeb(w, r)
}

func (h *Handlers) serveScheduleEventPage(w http.ResponseWriter, r *http.Request, event *models.ScheduledEvent) bool {
	page, err := fs.ReadFile(static.GetWeb(), "schedule/index.html")
	if err != nil {
		return false
	}

	baseURL := ""
	if h.configRepository != nil {
		baseURL = strings.TrimRight(h.configRepository.GetServerURL(), "/")
	}
	if baseURL == "" {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, r.Host)
	}
	pageURL := baseURL + r.URL.EscapedPath()
	imageURL := baseURL + "/logo/external"
	title := event.Name
	serverName := ""
	if h.configRepository != nil {
		serverName = h.configRepository.GetServerName()
	}
	if serverName != "" {
		title += " | " + serverName
	}
	description := event.Description
	if description == "" {
		description = title
	}
	endTime := event.StartTime.Add(time.Duration(event.DurationMinutes) * time.Minute)
	metadata := fmt.Sprintf(`<meta property="og:title" content="%s"/><meta property="og:description" content="%s"/><meta property="og:type" content="event"/><meta property="og:url" content="%s"/><meta property="og:site_name" content="%s"/><meta property="og:image" content="%s"/><meta property="og:image:alt" content="%s"/><meta property="event:start_time" content="%s"/><meta property="event:end_time" content="%s"/><meta name="twitter:card" content="summary"/><meta name="twitter:title" content="%s"/><meta name="twitter:description" content="%s"/><meta name="twitter:image" content="%s"/>`,
		html.EscapeString(title),
		html.EscapeString(description),
		html.EscapeString(pageURL),
		html.EscapeString(serverName),
		html.EscapeString(imageURL),
		html.EscapeString(title),
		html.EscapeString(event.StartTime.UTC().Format(time.RFC3339)),
		html.EscapeString(endTime.UTC().Format(time.RFC3339)),
		html.EscapeString(title),
		html.EscapeString(description),
		html.EscapeString(imageURL),
	)
	page = bytes.Replace(page, []byte("</head>"), []byte(metadata+"</head>"), 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(page)
	return true
}
