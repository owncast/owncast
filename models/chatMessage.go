package models

import (
	"bytes"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"mvdan.cc/xurls"
)

//ChatMessage represents a single chat message
type ChatMessage struct {
	ClientID string `json:"-"`

	Author      string    `json:"author"`
	Body        string    `json:"body"`
	Image       string    `json:"image"`
	ID          string    `json:"id"`
	MessageType string    `json:"type"`
	Visible     bool      `json:"visible"`
	Timestamp   time.Time `json:"timestamp"`
}

//Valid checks to ensure the message is valid
func (m ChatMessage) Valid() bool {
	return m.Author != "" && m.Body != "" && m.ID != ""
}

// RenderAndSanitizeMessageBody will turn markdown into HTML, sanitize raw user-supplied HTML and standardize
// the message into something safe and renderable for clients.
func (m *ChatMessage) RenderAndSanitizeMessageBody() {
	raw := m.Body

	// Set the new, sanitized and rendered message body
	m.Body = RenderAndSanitize(raw)
}

// RenderAndSanitize will turn markdown into HTML, sanitize raw user-supplied HTML and standardize
// the message into something safe and renderable for clients.
func RenderAndSanitize(raw string) string {
	rendered := renderMarkdown(raw)
	safe := sanitize(rendered)

	// Set the new, sanitized and rendered message body
	return strings.TrimSpace(safe)
}

func renderMarkdown(raw string) string {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			// html.WithXHTML(),
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

	// Allow links
	p.AllowAttrs("href").OnElements("a")

	// Force all URLs to have "noreferrer" in their rel attribute.
	p.RequireNoReferrerOnLinks(true)

	// Links will get target="_blank" added to them.
	p.AddTargetBlankToFullyQualifiedLinks(true)

	// Allow paragraphs
	p.AllowElements("br")
	p.AllowElements("p")

	// Allow img tags
	p.AllowElements("img")
	p.AllowAttrs("src").OnElements("img")
	p.AllowAttrs("alt").OnElements("img")
	p.AllowAttrs("title").OnElements("img")

	// Custom emoji have a class already specified.
	// We should only allow classes on emoji, not *all* imgs.
	// But TODO.
	p.AllowAttrs("class").OnElements("img")

	// Allow bold
	p.AllowElements("strong")

	// Allow emphasis
	p.AllowElements("em")

	return p.Sanitize(raw)
}
