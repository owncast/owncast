package plugins

import (
	"strings"
	"testing"
)

func loadedWithCommands(slug, display string, cmds []CommandInfo) *Loaded {
	return &Loaded{Manifest: &Manifest{Slug: slug, DisplayName: display, Commands: cmds}}
}

func TestIsHelpCommand(t *testing.T) {
	for _, b := range []string{"!help", "  !help  ", "!HELP", "!commands"} {
		if !IsHelpCommand(b) {
			t.Errorf("IsHelpCommand(%q) = false, want true", b)
		}
	}
	for _, b := range []string{"help", "!helpme", "!uptime", ""} {
		if IsHelpCommand(b) {
			t.Errorf("IsHelpCommand(%q) = true, want false", b)
		}
	}
}

func TestBuildHelpMessageEmpty(t *testing.T) {
	if got := BuildHelpMessage(nil, false); got != "" {
		t.Errorf("no plugins should yield empty help, got %q", got)
	}
	// A plugin with no commands also yields nothing.
	if got := BuildHelpMessage([]*Loaded{loadedWithCommands("x", "X", nil)}, false); got != "" {
		t.Errorf("plugin without commands should yield empty help, got %q", got)
	}
}

func TestBuildHelpMessageAggregatesAndGroups(t *testing.T) {
	plugins := []*Loaded{
		loadedWithCommands("stream-tracker", "Stream Tracker", []CommandInfo{
			{Name: "who", Prefix: "!", Description: "Who's in chat"},
			{Name: "uptime", Prefix: "!", Description: "Stream uptime"},
		}),
		loadedWithCommands("timer-bot", "Timer Bot", []CommandInfo{
			{Name: "remind", Prefix: "!", Description: "Set a reminder"},
		}),
	}
	out := BuildHelpMessage(plugins, false)
	for _, want := range []string{"Stream Tracker", "Timer Bot", "`!uptime`", "`!who`", "`!remind`", "Stream uptime"} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing %q\n---\n%s", want, out)
		}
	}
	// Grouped + sorted: Stream Tracker before Timer Bot; within it uptime before who.
	if strings.Index(out, "Stream Tracker") > strings.Index(out, "Timer Bot") {
		t.Error("plugins not sorted by display name")
	}
	if strings.Index(out, "`!uptime`") > strings.Index(out, "`!who`") {
		t.Error("commands not sorted within a plugin")
	}
}

func TestBuildHelpMessageModOnlyHiddenFromNonMods(t *testing.T) {
	plugins := []*Loaded{
		loadedWithCommands("mod", "Mod Tools", []CommandInfo{
			{Name: "ban", Prefix: "!", Description: "Ban a user", ModOnly: true},
			{Name: "rules", Prefix: "!", Description: "Show the rules"},
		}),
	}
	// Non-moderator: ban hidden, rules shown.
	nonMod := BuildHelpMessage(plugins, false)
	if strings.Contains(nonMod, "ban") {
		t.Errorf("mod-only command leaked to non-moderator:\n%s", nonMod)
	}
	if !strings.Contains(nonMod, "`!rules`") {
		t.Errorf("non-mod command missing for non-moderator:\n%s", nonMod)
	}
	// Moderator: both shown, ban tagged.
	mod := BuildHelpMessage(plugins, true)
	if !strings.Contains(mod, "`!ban`") || !strings.Contains(mod, "(mod)") {
		t.Errorf("moderator should see tagged mod-only command:\n%s", mod)
	}
}
