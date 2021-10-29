package events

import (
	"bytes"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/teris-io/shortid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"mvdan.cc/xurls"

	"github.com/owncast/owncast/core/user"
	log "github.com/sirupsen/logrus"
)

// EventPayload is a generic key/value map for sending out to chat clients.
type EventPayload map[string]interface{}

// OutboundEvent represents an event that is sent out to all listeners of the chat server.
type OutboundEvent interface {
	GetBroadcastPayload() EventPayload
	GetMessageType() EventType
}

// Event is any kind of event.  A type is required to be specified.
type Event struct {
	Type      EventType `json:"type,omitempty"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// UserEvent is an event with an associated user.
type UserEvent struct {
	User     *user.User `json:"user"`
	ClientID uint       `json:"clientId"`
	HiddenAt *time.Time `json:"hiddenAt,omitempty"`
}

// MessageEvent is an event that has a message body.
type MessageEvent struct {
	OutboundEvent `json:"-"`
	Body          string `json:"body"`
	RawBody       string `json:"-"`
}

// SystemActionEvent is an event that represents an action that took place, not a chat message.
type SystemActionEvent struct {
	Event
	MessageEvent
}

// SetDefaults will set default properties of all inbound events.
func (e *Event) SetDefaults() {
	e.ID = shortid.MustGenerate()
	e.Timestamp = time.Now()
}

// SetDefaults will set default properties of all inbound events.
func (e *UserMessageEvent) SetDefaults() {
	e.ID = shortid.MustGenerate()
	e.Timestamp = time.Now()
	e.RenderAndSanitizeMessageBody()
}

// RenderAndSanitizeMessageBody will turn markdown into HTML, sanitize raw user-supplied HTML and standardize
// the message into something safe and renderable for clients.
func (m *MessageEvent) RenderAndSanitizeMessageBody() {
	m.RawBody = m.Body

	// Set the new, sanitized and rendered message body
	m.Body = RenderAndSanitize(m.RawBody)
}

// Empty will return if this message's contents is empty.
func (m *MessageEvent) Empty() bool {
	return m.Body == ""
}

// RenderBody will render markdown to html without any sanitization.
func (m *MessageEvent) RenderBody() {
	m.RawBody = m.Body
	m.Body = RenderMarkdown(m.RawBody)
}

// RenderAndSanitize will turn markdown into HTML, sanitize raw user-supplied HTML and standardize
// the message into something safe and renderable for clients.
func RenderAndSanitize(raw string) string {
	rendered := RenderMarkdown(raw)
	safe := sanitize(rendered)

	// Set the new, sanitized and rendered message body
	return strings.TrimSpace(safe)
}

// RenderMarkdown will return HTML rendered from the string body of a chat message.
func RenderMarkdown(raw string) string {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{
					[]byte("http:"),
					[]byte("https:"),
				}),
				extension.WithLinkifyURLRegexp(
					xurls.Strict,
				),
			),
		),
	)

	trimmed := strings.TrimSpace(raw)
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(trimmed), &buf); err != nil {
		log.Debugln(err)
	}

	return buf.String()
}

var (
	_sanitizeReSrcMatch      = regexp.MustCompile(`(?i)^/img/emoji`)
	_sanitizeReAltTitleMatch = regexp.MustCompile(`:\S+:`)
)

func sanitize(raw string) string {
	p := bluemonday.StrictPolicy()

	// Require URLs to be parseable by net/url.Parse
	p.AllowStandardURLs()
	p.RequireParseableURLs(true)

	// Allow links
	p.AllowAttrs("href").OnElements("a")

	// Force all URLs to have "noreferrer" in their rel attribute.
	p.RequireNoReferrerOnLinks(true)

	// Links will get target="_blank" added to them.
	p.AddTargetBlankToFullyQualifiedLinks(true)

	// Allow breaks
	p.AllowElements("br")

	p.AllowElementsContent("p")

	// Allow img tags from the the local emoji directory only
	p.AllowAttrs("src").Matching(_sanitizeReSrcMatch).OnElements("img")
	p.AllowAttrs("alt", "title").Matching(_sanitizeReAltTitleMatch).OnElements("img")
	p.AllowAttrs("class").OnElements("img")

	// Allow bold
	p.AllowElements("strong")

	// Allow emphasis
	p.AllowElements("em")

	// Allow code blocks
	p.AllowElements("code")
	p.AllowElements("pre")

	return p.Sanitize(raw)
}
