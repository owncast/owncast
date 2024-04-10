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
