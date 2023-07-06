package events

import (
	"bytes"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/teris-io/shortid"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	emojiAst "github.com/yuin/goldmark-emoji/ast"
	emojiDef "github.com/yuin/goldmark-emoji/definition"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"mvdan.cc/xurls"

	"github.com/owncast/owncast/core/data"
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
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type,omitempty"`
	ID        string    `json:"id"`
}

// UserEvent is an event with an associated user.
type UserEvent struct {
	User     *user.User `json:"user"`
	HiddenAt *time.Time `json:"hiddenAt,omitempty"`
	ClientID uint       `json:"clientId,omitempty"`
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

// implements the emojiDef.Emojis interface but uses case-insensitive search.
// the .children field isn't currently used, but could be used in a future
// implementation of say, emoji packs where a child represents a pack.
type emojis struct {
	list     []emojiDef.Emoji
	names    map[string]*emojiDef.Emoji
	children []emojiDef.Emojis
}

// return a new Emojis set.
func newEmojis(emotes ...emojiDef.Emoji) emojiDef.Emojis {
	self := &emojis{
		list:     emotes,
		names:    map[string]*emojiDef.Emoji{},
		children: []emojiDef.Emojis{},
	}

	for i := range self.list {
		emoji := &self.list[i]
		for _, s := range emoji.ShortNames {
			self.names[s] = emoji
		}
	}

	return self
}

func (self *emojis) Get(shortName string) (*emojiDef.Emoji, bool) {
	v, ok := self.names[strings.ToLower(shortName)]
	if ok {
		return v, ok
	}

	for _, child := range self.children {
		v, ok := child.Get(shortName)
		if ok {
			return v, ok
		}
	}

	return nil, false
}

func (self *emojis) Add(emotes emojiDef.Emojis) {
	self.children = append(self.children, emotes)
}

func (self *emojis) Clone() emojiDef.Emojis {
	clone := &emojis{
		list:     self.list,
		names:    self.names,
		children: make([]emojiDef.Emojis, len(self.children)),
	}

	copy(clone.children, self.children)

	return clone
}

var (
	emojiMu           sync.Mutex
	emojiDefs         = newEmojis()
	emojiHTML         = make(map[string]string)
	emojiModTime      time.Time
	emojiHTMLFormat   = `<img src="{{ .URL }}" class="emoji" alt=":{{ .Name }}:" title=":{{ .Name }}:">`
	emojiHTMLTemplate = template.Must(template.New("emojiHTML").Parse(emojiHTMLFormat))
)

func loadEmoji() {
	modTime, err := data.UpdateEmojiList(false)
	if err != nil {
		return
	}

	if modTime.After(emojiModTime) {
		emojiMu.Lock()
		defer emojiMu.Unlock()

		emojiHTML = make(map[string]string)

		emojiList := data.GetEmojiList()
		emojiArr := make([]emojiDef.Emoji, 0)

		for i := 0; i < len(emojiList); i++ {
			var buf bytes.Buffer
			err := emojiHTMLTemplate.Execute(&buf, emojiList[i])
			if err != nil {
				return
			}
			emojiHTML[strings.ToLower(emojiList[i].Name)] = buf.String()

			emoji := emojiDef.NewEmoji(emojiList[i].Name, nil, strings.ToLower(emojiList[i].Name))
			emojiArr = append(emojiArr, emoji)
		}

		emojiDefs = newEmojis(emojiArr...)
	}
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
	loadEmoji()

	emojiMu.Lock()
	defer emojiMu.Unlock()

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
