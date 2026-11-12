package chat

import (
	"html"
	"regexp"
	"strings"

	goaway "github.com/TwiN/go-away"
)

// ChatMessageFilter is a allow/deny chat message filter.
type ChatMessageFilter struct{}

// NewMessageFilter will return an instance of the chat message filter.
func NewMessageFilter() *ChatMessageFilter {
	return &ChatMessageFilter{}
}

var (
	// Comments and tags are removed by the renderer before the message is
	// shown, so the filter must remove them too, or it checks a different text
	// than the one a viewer reads.
	htmlCommentPattern = regexp.MustCompile(`(?s)<!--.*?-->`)
	htmlTagPattern     = regexp.MustCompile(`<[^>]*>`)
	// Zero-width characters are invisible once rendered, so removing them
	// restores the word the reader sees. A non-breaking space is visible as a
	// space, so it becomes one.
	invisiblePattern = strings.NewReplacer(
		"\u200b", "",
		"\u200c", "",
		"\u200d", "",
		"\ufeff", "",
		"\u00a0", " ",
	)
)

// normalizeMessage renders a message the way a viewer will read it: entities
// resolved, tags and comments gone, invisible characters removed. The filter
// runs before the message is rendered and sanitized, so without this step a
// slur spelled with an entity, split by a tag or hidden behind a zero-width
// space still reaches the chat.
func normalizeMessage(message string) string {
	message = html.UnescapeString(message)
	message = htmlCommentPattern.ReplaceAllString(message, "")
	message = htmlTagPattern.ReplaceAllString(message, "")
	message = invisiblePattern.Replace(message)
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
