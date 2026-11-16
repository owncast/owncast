package chat

import (
	"html"
	"regexp"
	"strings"
	"unicode"

	goaway "github.com/TwiN/go-away"

	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/utils"
)

// ChatMessageFilter is a allow/deny chat message filter.
type ChatMessageFilter struct{}

// NewMessageFilter will return an instance of the chat message filter.
func NewMessageFilter() *ChatMessageFilter {
	return &ChatMessageFilter{}
}

var (
	// Comments are removed by the renderer before the message is shown, so
	// the filter must remove them too.
	htmlCommentPattern = regexp.MustCompile(`(?s)<!--.*?-->`)
)

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

// normalizeMessage renders a message the way a viewer will read it: entities
// resolved, tags and comments gone, invisible characters removed. The filter
// runs before the message is rendered and sanitized, so without this step a
// slur spelled with an entity, split by a tag or hidden behind a zero-width
// space still reaches the chat.
func normalizeMessage(message string) string {
	// Render the message the way a viewer reads it, so markdown syntax that
	// produces text (links, emphasis, images, ...) is reflected in the
	// filtered string instead of the raw markup.
	message = events.RenderMarkdown(message)
	message = html.UnescapeString(message)
	message = htmlCommentPattern.ReplaceAllString(message, "")
	// bluemonday strips tags (including malformed ones) more reliably than a
	// regex, matching what the renderer's sanitizer leaves as text.
	message = utils.StripHTML(message)
	message = dropInvisible(message)
	return strings.Join(strings.Fields(message), " ")
}

// Allow will test if this message should be allowed to be sent.
func (*ChatMessageFilter) Allow(message string) bool {
	// The raw message is checked as well: normalization strips markup, and a
	// message that was profane before that strip still is.
	if goaway.IsProfane(message) {
		return false
	}
	return !goaway.IsProfane(normalizeMessage(message))
}
