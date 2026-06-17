package inbox

import (
	"strings"
	"testing"

	"github.com/owncast/owncast/services/activitypub/apmodels"
)

func TestTruncateMetadata(t *testing.T) {
	if got := truncateMetadata("hello", 100); got != "hello" {
		t.Errorf("short string should be unchanged, got %q", got)
	}
	if got := truncateMetadata("hello world", 5); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
	// UTF-8 safety: truncation must not split a multi-byte rune.
	if got := truncateMetadata("héllo", 2); got != "hé" {
		t.Errorf("expected %q, got %q", "hé", got)
	}
	if got := truncateMetadata("abc", 0); got != "" {
		t.Errorf("max 0 should yield empty, got %q", got)
	}
}

func TestClampTags(t *testing.T) {
	many := make([]string, maxTags+10)
	for i := range many {
		many[i] = "tag"
	}
	if got := clampTags(many); len(got) != maxTags {
		t.Errorf("expected tag count clamped to %d, got %d", maxTags, len(got))
	}

	long := []string{strings.Repeat("x", maxTagLen+50)}
	clamped := clampTags(long)
	if len([]rune(clamped[0])) != maxTagLen {
		t.Errorf("expected each tag clamped to %d runes, got %d", maxTagLen, len([]rune(clamped[0])))
	}
}

func TestBuildStreamUpdateClampsHostileMetadata(t *testing.T) {
	hostile := &apmodels.OwncastMetadata{
		StreamTitle:       strings.Repeat("T", maxStreamTitleLen+500),
		StreamDescription: strings.Repeat("D", maxStreamDescriptionLen+500),
		ThumbnailURL:      "https://example.com/" + strings.Repeat("a", maxMetadataURLLen+500),
		Tags:              make([]string, maxTags+5),
	}

	update := buildStreamUpdateFromMetadata(hostile)

	if update.Title == nil || len([]rune(*update.Title)) != maxStreamTitleLen {
		t.Errorf("stream title not clamped to %d", maxStreamTitleLen)
	}
	if update.Description == nil || len([]rune(*update.Description)) != maxStreamDescriptionLen {
		t.Errorf("stream description not clamped to %d", maxStreamDescriptionLen)
	}
	if update.ThumbnailURL == nil || len([]rune(*update.ThumbnailURL)) != maxMetadataURLLen {
		t.Errorf("thumbnail URL not clamped to %d", maxMetadataURLLen)
	}
	if len(update.Tags) != maxTags {
		t.Errorf("tags not clamped to %d, got %d", maxTags, len(update.Tags))
	}
}
