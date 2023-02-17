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

	if len(newString) > length {
		newString = newString[:length]
	}

	newString = strings.ReplaceAll(newString, "\r", "")
	newString = strings.TrimSpace(newString)

	return newString
}
