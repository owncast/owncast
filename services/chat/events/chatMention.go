package events

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/owncast/owncast/models"
	htmlparser "golang.org/x/net/html"
)

// ChatMention identifies a user explicitly selected in the composer. The
// offsets refer to UTF-8 bytes in the original Markdown body.
type ChatMention struct {
	UserID string `json:"userId"`
	Start  int    `json:"start"`
	End    int    `json:"end"`
}

// RenderChatMentions preserves selected labels across Markdown rendering. Only
// server-validated text nodes can become mention markup in the stored body.
func RenderChatMentions(raw, fallback string, mentions []ChatMention, getUser func(string) *models.User) string {
	const maxMentions = 16
	if len(mentions) == 0 || len(mentions) > maxMentions {
		return fallback
	}

	prefix := "owncastmentiontoken"
	for strings.Contains(raw, prefix) {
		prefix += "x"
	}

	var markdown strings.Builder
	replacements := make([]string, 0, 2*len(mentions))
	lastEnd := 0
	for _, mention := range mentions {
		if mention.Start < lastEnd || mention.End <= mention.Start || mention.End > len(raw) {
			continue
		}
		label := raw[mention.Start:mention.End]
		if !strings.HasPrefix(label, "@") {
			continue
		}
		user := getUser(mention.UserID)
		if user == nil || user.DisabledAt != nil || user.IsBot || label != "@"+user.DisplayName {
			continue
		}

		marker := fmt.Sprintf("%s%dz", prefix, len(replacements)/2)
		// Include the selected identity in the stored HTML. The normal sanitizer
		// disallows span, so only validated selections can produce these tags.
		markup := `<span data-chat-mention="` + html.EscapeString(user.ID) + `">` + html.EscapeString(label) + "</span>"
		markdown.WriteString(raw[lastEnd:mention.Start])
		markdown.WriteString(marker)
		replacements = append(replacements, marker, markup)
		lastEnd = mention.End
	}

	if lastEnd == 0 {
		return fallback
	}
	markdown.WriteString(raw[lastEnd:])

	return replaceChatMentionTokens(RenderAndSanitize(markdown.String()), fallback, prefix, replacements)
}

func replaceChatMentionTokens(rendered, fallback, prefix string, replacements []string) string {
	replace := strings.NewReplacer(replacements...)
	tokenizer := htmlparser.NewTokenizer(strings.NewReader(rendered))
	var output bytes.Buffer
	count := 0
	for kind := tokenizer.Next(); kind != htmlparser.ErrorToken; kind = tokenizer.Next() {
		token := string(tokenizer.Raw())
		found := strings.Count(token, prefix)
		if found > 0 {
			if kind != htmlparser.TextToken {
				return fallback
			}
			token = replace.Replace(token)
			count += found
		}
		output.WriteString(token)
	}
	if count != len(replacements)/2 {
		return fallback
	}
	return output.String()
}
