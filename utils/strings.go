package utils

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// StripHTML will strip HTML tags from a string.
func StripHTML(s string) string {
	p := bluemonday.NewPolicy()
	return p.Sanitize(s)
}

// MakeSafeStringOfLength will take a string and strip HTML tags,
// trim whitespace, and limit the length.
func MakeSafeStringOfLength(s string, length int) string {
	newString, fullyUnescaped := unescapeHTML(s)
	newString = StripHTML(newString)
	// Only decode display entities after the input is stable. If the pass limit
	// was reached, preserving the encoding prevents latent markup from activating.
	if fullyUnescaped {
		newString = html.UnescapeString(newString)
	}

	// Convert utf-8 string into Unicode code points.
	codePoints := []rune(newString)

	if len(codePoints) > length {
		codePoints = codePoints[:length]
	}

	newString = string(codePoints)
	newString = strings.ReplaceAll(newString, "\r", "")
	newString = strings.TrimSpace(newString)

	return newString
}

func unescapeHTML(s string) (string, bool) {
	const maxPasses = 8
	for range maxPasses {
		unescaped := html.UnescapeString(s)
		if unescaped == s {
			return s, true
		}
		s = unescaped
	}
	return s, false
}
