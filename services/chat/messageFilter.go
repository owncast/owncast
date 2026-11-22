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

// altTextPattern captures the value of an alt attribute on an image. It is used
// to build the fallback text a viewer sees when an image fails to load.
var altTextPattern = regexp.MustCompile(`(?i)[[:space:]]alt[[:space:]]*=[[:space:]]*(?:"([^"]*)"|'([^']*)')`)

// imgTagPattern matches a whole image tag so it can be replaced in place by its
// alt text.
var imgTagPattern = regexp.MustCompile(`(?i)<img\b[^>]*>`)

// normalizeFilterText unescapes entities, drops invisible characters and
// collapses whitespace so the checked text matches what a viewer reads.
func normalizeFilterText(s string) string {
	s = html.UnescapeString(s)
	s = dropInvisible(s)
	return strings.Join(strings.Fields(s), " ")
}

// AllowRendered checks the rendered message body (HTML from SetDefaults) for
// profanity. Two representations are checked independently:
//
//   - the image-free body, which is what a viewer reads when images load;
//   - a positional fallback where each image is replaced by its alt text at
//     that image's position, which is what a viewer reads when images fail.
//
// Checking each representation on its own matters. Appending every alt value
// to the end of the stripped body can join a profane body word with a benign
// alt word into an allowlisted token (e.g. "cock" + "tail" = "cocktail"), and
// it can reorder text so profanity split around an image is no longer seen
// (e.g. "f" + "uc" + "k" becomes "fk uc"). Neither representation may be
// derived from the other, or the filter can be bypassed.
func (*ChatMessageFilter) AllowRendered(renderedBody string) bool {
	// Image-free body: StripHTML removes img tags entirely, including alt.
	imageFree := normalizeFilterText(utils.StripHTML(renderedBody))
	if goaway.IsProfane(imageFree) {
		return false
	}

	// Positional fallback: replace each image with its alt text where it sits
	// in the document, then check the resulting text.
	fallback := imgTagPattern.ReplaceAllStringFunc(renderedBody, func(tag string) string {
		m := altTextPattern.FindStringSubmatch(tag)
		if m == nil {
			return ""
		}
		if m[1] != "" {
			return m[1]
		}
		return m[2]
	})
	fallback = normalizeFilterText(utils.StripHTML(fallback))
	return !goaway.IsProfane(fallback)
}
