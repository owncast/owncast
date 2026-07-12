package plugins

import (
	"math"
	"reflect"
	"testing"
)

func TestMatchCommand(t *testing.T) {
	msg := HostChatMessage{Body: "?PING one  two"}
	command := CommandInfo{Name: "ping", Prefix: "?", Aliases: []string{"p"}}

	event, matched := matchCommand(msg, command)
	if !matched {
		t.Fatal("case-insensitive command should match")
	}
	if event.Command != "ping" || event.InvokedAs != "PING" {
		t.Fatalf("unexpected command identity: %+v", event)
	}
	if !reflect.DeepEqual(event.Args, []string{"one", "two"}) || event.ArgString != "one  two" {
		t.Fatalf("unexpected arguments: %+v", event)
	}

	command.CaseSensitive = true
	if _, matched := matchCommand(msg, command); matched {
		t.Fatal("case-sensitive command should reject different casing")
	}
	if _, matched := matchCommand(HostChatMessage{Body: "?unknown"}, command); matched {
		t.Fatal("unknown command should not match")
	}
}

func TestCommandTargetsAllowDuplicatesAndIgnoreUnknown(t *testing.T) {
	first := &Loaded{Manifest: &Manifest{Slug: "first", Commands: []CommandInfo{{Name: "ping"}}}}
	second := &Loaded{Manifest: &Manifest{Slug: "second", Commands: []CommandInfo{{Name: "ping"}}}}
	dispatcher := NewDispatcher([]*Loaded{first, second})

	targets := dispatcher.commandTargets(HostChatMessage{Body: "!ping"})
	if len(targets) != 2 {
		t.Fatalf("duplicate declarations should both match, got %d targets", len(targets))
	}
	if targets[0].plugin != first || targets[1].plugin != second {
		t.Fatalf("unexpected command targets: %+v", targets)
	}
	if targets := dispatcher.commandTargets(HostChatMessage{Body: "!unknown"}); len(targets) != 0 {
		t.Fatalf("unknown command should be silent, got %d targets", len(targets))
	}
	first.Manifest.Commands = append(first.Manifest.Commands, CommandInfo{Name: "help"})
	if targets := dispatcher.commandTargets(HostChatMessage{Body: HelpCommand}); len(targets) != 1 || targets[0].plugin != first {
		t.Fatalf("built-in help command should remain available to plugins: %+v", targets)
	}
}

func TestCommandModeratorGate(t *testing.T) {
	command := CommandInfo{Name: "announce", ModOnly: true}
	if commandAllowed(HostChatMessage{}, command) {
		t.Fatal("anonymous sender must not pass moderator gate")
	}
	viewer := HostChatMessage{User: &HostUser{ID: "viewer"}}
	if commandAllowed(viewer, command) {
		t.Fatal("viewer must not pass moderator gate")
	}
	moderator := HostChatMessage{User: &HostUser{ID: "mod", Scopes: []string{commandModeratorScope}}}
	if !commandAllowed(moderator, command) {
		t.Fatal("moderator should pass moderator gate")
	}
}

func TestCommandCooldownIsolation(t *testing.T) {
	loaded := &Loaded{}
	command := CommandInfo{Name: "ping", CooldownMs: 1_000}
	message := func(user, timestamp string) HostChatMessage {
		return HostChatMessage{User: &HostUser{ID: user}, Timestamp: timestamp}
	}

	if !commandCooldownReady(loaded, message("u1", "2026-01-01T00:00:00Z"), command) {
		t.Fatal("first invocation should run")
	}
	if commandCooldownReady(loaded, message("u1", "2026-01-01T00:00:00.500Z"), command) {
		t.Fatal("same user and command should be cooling down")
	}
	if !commandCooldownReady(loaded, message("u2", "2026-01-01T00:00:00.500Z"), command) {
		t.Fatal("cooldown must be isolated by user")
	}
	other := CommandInfo{Name: "other", CooldownMs: 1_000}
	if !commandCooldownReady(loaded, message("u1", "2026-01-01T00:00:00.500Z"), other) {
		t.Fatal("cooldown must be isolated by command")
	}
	if !commandCooldownReady(loaded, message("u1", "2026-01-01T00:00:01Z"), command) {
		t.Fatal("invocation should run when cooldown expires")
	}

	if !commandCooldownReady(loaded, message("u3", "2026-01-01T00:02:00Z"), command) {
		t.Fatal("later invocation should run")
	}
	if len(loaded.commandCooldowns) != 1 {
		t.Fatalf("expired cooldown entries were not pruned: %+v", loaded.commandCooldowns)
	}
}

func TestCommandCooldownOverflow(t *testing.T) {
	loaded := &Loaded{}
	command := CommandInfo{Name: "ping", CooldownMs: math.MaxInt64}
	first := HostChatMessage{User: &HostUser{ID: "u1"}, Timestamp: "2026-01-01T00:00:00Z"}
	if !commandCooldownReady(loaded, first, command) {
		t.Fatal("first invocation should run")
	}
	second := HostChatMessage{User: &HostUser{ID: "u1"}, Timestamp: "2026-01-01T00:00:01Z"}
	if commandCooldownReady(loaded, second, command) {
		t.Fatal("overflowing cooldown must remain active")
	}
}
