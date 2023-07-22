package chat

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/owncast/owncast/models"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	emojiAst "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"mvdan.cc/xurls"

	log "github.com/sirupsen/logrus"
)

// OutboundEvent represents an event that is sent out to all listeners of the chat server.
type OutboundEvent interface {
	GetBroadcastPayload() models.EventPayload
	GetMessageType() models.EventType
}

// MessageEvent is an event that has a message body.
type MessageEvent struct {
	OutboundEvent `json:"-"`
	Body          string `json:"body"`
	RawBody       string `json:"-"`
}

// SystemActionEvent is an event that represents an action that took place, not a chat message.
type SystemActionEvent struct {
	models.Event
	MessageEvent
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
	// emojiMu.Lock()
	// defer emojiMu.Unlock()

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
			emoji.New(
				emoji.WithEmojis(
					emojiDefs,
				),
				emoji.WithRenderingMethod(emoji.Func),
				emoji.WithRendererFunc(func(w util.BufWriter, source []byte, n *emojiAst.Emoji, config *emoji.RendererConfig) {
					baseName := n.Value.ShortNames[0]
					_, _ = w.WriteString(emojiHTML[baseName])
				}),
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
	_sanitizeReSrcMatch    = regexp.MustCompile(`(?i)^/img/emoji/[^\.%]*.[A-Z]*$`)
	_sanitizeReClassMatch  = regexp.MustCompile(`(?i)^(emoji)[A-Z_]*?$`)
	_sanitizeNonEmptyMatch = regexp.MustCompile(`^.+$`)
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

	p.AllowElements("p")

	// Allow img tags from the the local emoji directory only
	p.AllowAttrs("src").Matching(_sanitizeReSrcMatch).OnElements("img")
	p.AllowAttrs("alt", "title").Matching(_sanitizeNonEmptyMatch).OnElements("img")
	p.AllowAttrs("class").Matching(_sanitizeReClassMatch).OnElements("img")

	// Allow bold
	p.AllowElements("strong")

	// Allow emphasis
	p.AllowElements("em")

	// Allow code blocks
	p.AllowElements("code")
	p.AllowElements("pre")

	return p.Sanitize(raw)
}
