package chat

import (
	"testing"
)

func TestFiltering(t *testing.T) {
	filter := NewMessageFilter()

	filteredTestMessages := []string{
		"Hello, fucking world!",
		"Suck my dick",
		"Eat my ass",
		"fuck this shit",
		"@$$h073",
		"F   u   C  k th1$ $h!t",
		"u r fag",
		"fucking sucks",
	}

	unfilteredTestMessages := []string{
		"bass fish",
		"assumptions",
	}

	for _, m := range filteredTestMessages {
		result := filter.Allow(m)
		if result {
			t.Errorf("%s should be seen as a filtered profane message", m)
		}
	}

	for _, m := range unfilteredTestMessages {
		result := filter.Allow(m)
		if !result {
			t.Errorf("%s should not be filtered", m)
		}
	}
}

// The filter runs on the raw payload, while the viewer reads the rendered text:
// entities are resolved and tags, comments and invisible characters are gone by
// then. These messages read as profane and used to pass.
func TestFilteringSeesTheRenderedText(t *testing.T) {
	filter := NewMessageFilter()

	filteredTestMessages := []string{
		"&#102;uck this",           // entity resolves to f in the rendered text
		"fu<b>ck</b> this",         // a tag is stripped, the letters join
		"f<!-- hidden -->uck this", // a comment is stripped, the letters join
		"fu\u200bck this",          // zero-width space is invisible
		"fu\u00a0ck this",          // non-breaking space reads as a space
		"<span>fuck</span> this",   // tag-wrapped, plainly profane
	}

	unfilteredTestMessages := []string{
		"<b>bold</b> text",
		"class &amp; assignment",
		"nice work &#128077;",
		"assumptions about the bass fish",
	}

	for _, m := range filteredTestMessages {
		if filter.Allow(m) {
			t.Errorf("%q should be filtered: the rendered text is profane", m)
		}
	}

	for _, m := range unfilteredTestMessages {
		if !filter.Allow(m) {
			t.Errorf("%q should not be filtered", m)
		}
	}
}
