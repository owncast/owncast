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
		result := filter.AllowRaw(m)
		if result {
			t.Errorf("%s should be seen as a filtered profane message", m)
		}
	}

	for _, m := range unfilteredTestMessages {
		result := filter.AllowRaw(m)
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

	// These are HTML-rendered bodies (as produced by SetDefaults).
	// They contain profanity once tags are stripped and entities unescaped.
	filteredTestMessages := []string{
		"<p>&#102;uck this</p>",                                         // entity resolves to f in the rendered text
		"<p>fu<b>ck</b> this</p>",                                       // a tag is stripped, the letters join
		"<p>f<!-- hidden -->uck this</p>",                               // a comment is stripped, the letters join
		"<p>fu\u200bck this</p>",                                        // zero-width space is invisible
		"<p>fu\u00a0ck this</p>",                                        // non-breaking space reads as a space
		"<p><span>fuck</span> this</p>",                                 // tag-wrapped, plainly profane
		`<p><img src="/img/emoji/nonexistent.svg" alt="&#102;uck"></p>`, // alt text shown on failed image load
		"<p>f<a href=\"http://x\">u</a>ck this</p>",                     // markdown link becomes <a>, text joins
		"<p>fu<strong>ck</strong> this</p>",                             // markdown emphasis becomes <strong>, letters join
	}

	unfilteredTestMessages := []string{
		"<p><b>bold</b> text</p>",
		"<p>class & assignment</p>",
		"<p>nice work &#128077;</p>",
		"<p>assumptions about the bass fish</p>",
	}

	for _, m := range filteredTestMessages {
		if filter.AllowRendered(m) {
			t.Errorf("%q should be filtered: the rendered text is profane", m)
		}
	}

	for _, m := range unfilteredTestMessages {
		if !filter.AllowRendered(m) {
			t.Errorf("%q should not be filtered", m)
		}
	}
}
