package events

import (
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
)

func TestRenderChatMentions(t *testing.T) {
	user := &models.User{ID: "alice", DisplayName: "Élodie & Co"}
	lookup := func(id string) *models.User {
		if id == user.ID {
			return user
		}
		return nil
	}
	label := "@" + user.DisplayName
	prefix := "Hi 👋 **"
	raw := prefix + label + "**, hello " + label
	mention := ChatMention{UserID: user.ID, Start: len(prefix), End: len(prefix + label)}
	result := RenderChatMentions(raw, RenderAndSanitize(raw), []ChatMention{mention}, lookup)
	if want := `<strong><span data-chat-mention="alice">@Élodie &amp; Co</span></strong>`; !strings.Contains(result, want) {
		t.Fatalf("rendered mention = %q, want it to contain %q", result, want)
	}
	if strings.Count(result, "<span ") != 1 || !strings.Contains(result, ", hello @Élodie &amp; Co") {
		t.Fatalf("only the selected occurrence should be highlighted: %q", result)
	}

	t.Run("invalid identities or offsets leave body unchanged", func(t *testing.T) {
		fallback := RenderAndSanitize(raw)
		for _, invalid := range []ChatMention{
			{UserID: "missing", Start: mention.Start, End: mention.End},
			{UserID: user.ID, Start: -1, End: mention.End},
			{UserID: user.ID, Start: mention.Start, End: mention.End - 1},
		} {
			if got := RenderChatMentions(raw, fallback, []ChatMention{invalid}, lookup); got != fallback {
				t.Errorf("invalid mention %+v changed the rendered body: %q", invalid, got)
			}
		}
	})

	t.Run("markup in a display name remains text", func(t *testing.T) {
		user.DisplayName = `<img src=x onerror="alert(1)">`
		raw := "@" + user.DisplayName
		got := RenderChatMentions(raw, RenderAndSanitize(raw), []ChatMention{{UserID: user.ID, End: len(raw)}}, lookup)
		if strings.Contains(got, "<img") || !strings.Contains(got, "@&lt;img") {
			t.Fatalf("display name must be escaped: %q", got)
		}
	})

	t.Run("a mention in an attribute cannot inject HTML", func(t *testing.T) {
		user.DisplayName = "Alice"
		prefix := `<a href="https://example.com/`
		raw := prefix + `@Alice">example</a>`
		fallback := RenderAndSanitize(raw)
		mention := ChatMention{UserID: user.ID, Start: len(prefix), End: len(prefix + "@Alice")}
		if got := RenderChatMentions(raw, fallback, []ChatMention{mention}, lookup); got != fallback {
			t.Fatalf("mention inside attribute should keep original output: %q", got)
		}
	})
}
