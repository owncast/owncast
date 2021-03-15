package models

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
)

// ChatEvent represents a single chat message.
type ChatEvent struct {
	ClientID string `json:"-"`

	Author      string    `json:"author,omitempty"`
	Body        string    `json:"body,omitempty"`
	RawBody     string    `json:"-"`
	ID          string    `json:"id"`
	MessageType EventType `json:"type"`
	Visible     bool      `json:"visible"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Ephemeral   bool      `json:"ephemeral,omitempty"`
}

// Valid checks to ensure the message is valid.
func (m ChatEvent) Valid() bool {
	return m.Author != "" && m.Body != "" && m.ID != ""
}

// SetDefaults will set default values on a chat event object.
func (m *ChatEvent) SetDefaults() {
	id, _ := shortid.Generate()
	m.ID = id
	m.Timestamp = time.Now()
	m.Visible = true
}

// RenderAndSanitizeMessageBody will turn markdown into HTML, sanitize raw user-supplied HTML and standardize
// the message into something safe and renderable for clients.
func (m *ChatEvent) RenderAndSanitizeMessageBody() {
	m.RawBody = m.Body

	// Set the new, sanitized and rendered message body
	m.Body = RenderAndSanitize(m.RawBody)
}

// Empty will return if this message's contents is empty.
func (m *ChatEvent) Empty() bool {
	return m.Body == ""
}

// RenderBody will render markdown to html without any sanitization
func (m *ChatEvent) RenderBody() {
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
		panic(err)
	}

	return buf.String()
}

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
	p.AllowElements("br", "p")

	// Allow img tags from the the local emoji directory only
	p.AllowAttrs("src", "alt", "class", "title").Matching(regexp.MustCompile(`(?i)/img/emoji`)).OnElements("img")
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
