package chat

import (
	"testing"

	"github.com/owncast/owncast/core/chat/events"
)

// Test a bunch of arbitrary markup and markdown to make sure we get sanitized
// and fully rendered HTML out of it.
func TestRenderAndSanitize(t *testing.T) {
	messageContent := `
  Test one two three!  I go to http://yahoo.com and search for _sports_ and **answers**.
  Here is an iframe <iframe src="http://yahoo.com"></iframe>

  ## blah blah blah
  [test link](http://owncast.online)
  <img class="emoji" alt="bananadance.gif" width="600px" src="/img/emoji/bananadance.gif">
  <script src="http://hackers.org/hack.js"></script>
  `

	expected := `Test one two three!  I go to <a href="http://yahoo.com" rel="nofollow noreferrer noopener" target="_blank">http://yahoo.com</a> and search for <em>sports</em> and <strong>answers</strong>.
Here is an iframe 
blah blah blah
<a href="http://owncast.online" rel="nofollow noreferrer noopener" target="_blank">test link</a>
<img class="emoji" src="/img/emoji/bananadance.gif">`

	result := events.RenderAndSanitize(messageContent)
	if result != expected {
		t.Errorf("message rendering/sanitation does not match expected.  Got\n%s, \n\n want:\n%s", result, expected)
	}
}

// Test to make sure we block remote images in chat messages.
func TestBlockRemoteImages(t *testing.T) {
	messageContent := `<img src="https://via.placeholder.com/img/emoji/350x150"> test ![](https://via.placeholder.com/img/emoji/350x150)`
	expected := `test`
	result := events.RenderAndSanitize(messageContent)

	if result != expected {
		t.Errorf("message rendering/sanitation does not match expected.  Got\n%s, \n\n want:\n%s", result, expected)
	}
}

// Test to make sure emoji images are allowed in chat messages.
func TestAllowEmojiImages(t *testing.T) {
	messageContent := `<img alt=":beerparrot:" title=":beerparrot:" src="/img/emoji/beerparrot.gif"> test ![](/img/emoji/beerparrot.gif)`
	expected := `<img alt=":beerparrot:" title=":beerparrot:" src="/img/emoji/beerparrot.gif"> test <img src="/img/emoji/beerparrot.gif">`
	result := events.RenderAndSanitize(messageContent)

	if result != expected {
		t.Errorf("message rendering/sanitation does not match expected.  Got\n%s, \n\n want:\n%s", result, expected)
	}
}

// Test to verify we can pass raw html and render markdown.
func TestAllowHTML(t *testing.T) {
	messageContent := `<img src="/img/emoji/beerparrot.gif"><ul><li>**test thing**</li></ul>`
	expected := "<p><img src=\"/img/emoji/beerparrot.gif\"><ul><li><strong>test thing</strong></li></ul></p>\n"
	result := events.RenderMarkdown(messageContent)

	if result != expected {
		t.Errorf("message rendering does not match expected.  Got\n%s, \n\n want:\n%s", result, expected)
	}
}
