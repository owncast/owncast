package utils

import (
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
	newString := s
	newString = StripHTML(newString)

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
