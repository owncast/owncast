package plugins

import (
	"sort"
	"strings"
)

// HelpCommand is the chat command that triggers the unified command listing.
// The host owns it (no plugin needs to implement it); a plugin can't shadow it.
const HelpCommand = "!help"

// IsHelpCommand reports whether a chat message body is the help command. Matched
// case-insensitively and trim-tolerant; "!commands" is accepted as an alias.
func IsHelpCommand(body string) bool {
	b := strings.ToLower(strings.TrimSpace(body))
	return b == HelpCommand || b == "!commands"
}

// BuildHelpMessage renders the unified command list aggregated across the given
// loaded plugins, for a viewer whose moderator status is isModerator. Commands
// are grouped by plugin (so the source is clear) and sorted; mod-only commands
// are hidden from non-moderators. Returns "" when there's nothing to show, so
// the caller can stay silent rather than post an empty list.
func BuildHelpMessage(loaded []*Loaded, isModerator bool) string {
	type group struct {
		plugin string
		lines  []string
	}
	var groups []group

	// Stable plugin order by display name (fallback slug).
	sorted := append([]*Loaded(nil), loaded...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return helpPluginName(sorted[i]) < helpPluginName(sorted[j])
	})

	for _, l := range sorted {
		if l == nil || l.Manifest == nil {
			continue
		}
		cmds := append([]CommandInfo(nil), l.Manifest.Commands...)
		sort.SliceStable(cmds, func(i, j int) bool { return cmds[i].Name < cmds[j].Name })

		var lines []string
		for _, c := range cmds {
			if c.ModOnly && !isModerator {
				continue
			}
			prefix := c.Prefix
			if prefix == "" {
				prefix = "!"
			}
			// Markdown: backticked command, em-dash description. System
			// messages are rendered (RenderBody), so this becomes a code span
			// + text; it also degrades to readable plaintext if unrendered.
			line := "- `" + prefix + c.Name + "`"
			if c.Description != "" {
				line += " — " + c.Description
			}
			if c.ModOnly {
				line += " _(mod)_"
			}
			lines = append(lines, line)
		}
		if len(lines) > 0 {
			groups = append(groups, group{plugin: helpPluginName(l), lines: lines})
		}
	}

	if len(groups) == 0 {
		return ""
	}

	// Rendered as markdown: a header paragraph, then per-plugin a bold heading
	// followed by a bulleted list (blank lines separate the blocks so goldmark
	// renders the list, not a run-on line).
	var b strings.Builder
	b.WriteString("Available chat commands:")
	for _, g := range groups {
		b.WriteString("\n\n**")
		b.WriteString(g.plugin)
		b.WriteString("**\n\n")
		b.WriteString(strings.Join(g.lines, "\n"))
	}
	return b.String()
}

func helpPluginName(l *Loaded) string {
	if l == nil || l.Manifest == nil {
		return ""
	}
	if l.Manifest.DisplayName != "" {
		return l.Manifest.DisplayName
	}
	return l.Manifest.Slug
}
