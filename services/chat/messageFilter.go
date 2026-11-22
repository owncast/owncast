package chat

import (
	"html"
	"regexp"
	"strings"
	"unicode"

	goaway "github.com/TwiN/go-away"

	"github.com/owncast/owncast/utils"
)

// ChatMessageFilter is a allow/deny chat message filter.
type ChatMessageFilter struct{}

// NewMessageFilter will return an instance of the chat message filter.
func NewMessageFilter() *ChatMessageFilter {
	return &ChatMessageFilter{}
}

// dropInvisible removes format characters that are invisible once rendered
// (zero-width spaces, soft hyphens, word joiners, BOM, ...) and maps a
// non-breaking space to a regular space, so the filtered text matches what a
// viewer reads.
func dropInvisible(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\u00a0' {
			return ' '
		}
		if unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, s)
}

// AllowRaw checks the raw message payload for profanity.
// This is called early in readPump before the message is rendered.
func (*ChatMessageFilter) AllowRaw(message string) bool {
	return !goaway.IsProfane(message)
}

// AllowRendered checks the rendered message body (HTML from SetDefaults)
// for profanity. It strips tags, unescapes entities, and removes invisible
// characters to match what a viewer reads.
// altTextPattern captures the value of an alt attribute so the text a viewer
// can see when an image fails to load is also checked for profanity.
var altTextPattern = regexp.MustCompile(`(?i)[[:space:]]alt[[:space:]]*=[[:space:]]*(?:"([^"]*)"|'([^']*)')`)

func (*ChatMessageFilter) AllowRendered(renderedBody string) bool {
	// Collect alt text from images before stripping tags: a viewer whose image
	// fails to load can see the alt text, so it must be part of the checked
	// string.
	extra := ""
	for _, m := range altTextPattern.FindAllStringSubmatch(renderedBody, -1) {
		if m[1] != "" {
			extra += " " + m[1]
		} else {
			extra += " " + m[2]
		}
	}
	// Strip all HTML tags (bluemonday with no allowed elements).
	text := utils.StripHTML(renderedBody) + extra
	// Unescape HTML entities.
	text = html.UnescapeString(text)
	// Remove invisible format characters and normalize spaces.
	text = dropInvisible(text)
	// Collapse whitespace.
	text = strings.Join(strings.Fields(text), " ")
	return !goaway.IsProfane(text)
}
