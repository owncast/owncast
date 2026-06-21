package plugins

import (
	"context"
	"strings"
	"testing"
)

// seedManager builds a Manager with the given discovered entries and enabled
// set, bypassing disk scanning so we can unit-test the auth.gate guard and
// ActiveAuthGate without loading real wasm.
func seedManager(discovered map[string]*DiscoveredEntry, enabled map[string]bool, loaded map[string]*Loaded) *Manager {
	m := NewManager(".", &HostEnv{})
	m.discovered = discovered
	m.enabledSet = enabled
	if loaded != nil {
		m.loaded = loaded
	}
	return m
}

func TestManager_RefusesSecondAuthGate(t *testing.T) {
	m := seedManager(
		map[string]*DiscoveredEntry{
			"gate-a": {Slug: "gate-a", Permissions: []string{PermAuthGate}},
			"gate-b": {Slug: "gate-b", Permissions: []string{PermAuthGate}},
		},
		map[string]bool{"gate-a": true},
		nil,
	)

	err := m.Enable(context.Background(), "gate-b")
	if err == nil {
		t.Fatal("expected enabling a second auth.gate plugin to fail")
	}
	if !strings.Contains(err.Error(), "another authentication-gate plugin") {
		t.Fatalf("unexpected error: %v", err)
	}
	// The second gate must not have been added to the enabled set.
	if m.enabledSet["gate-b"] {
		t.Fatal("gate-b should not be enabled after the guard rejected it")
	}
}

func TestManager_ActiveAuthGate(t *testing.T) {
	// Enabled + loaded gate.
	m := seedManager(
		map[string]*DiscoveredEntry{"gate-a": {Slug: "gate-a", Permissions: []string{PermAuthGate}}},
		map[string]bool{"gate-a": true},
		map[string]*Loaded{"gate-a": {}},
	)
	if slug, loaded := m.ActiveAuthGate(); slug != "gate-a" || !loaded {
		t.Fatalf("ActiveAuthGate: got (%q,%v) want (gate-a,true)", slug, loaded)
	}

	// Enabled but not loaded → armed but unavailable.
	m2 := seedManager(
		map[string]*DiscoveredEntry{"gate-a": {Slug: "gate-a", Permissions: []string{PermAuthGate}}},
		map[string]bool{"gate-a": true},
		map[string]*Loaded{},
	)
	if slug, loaded := m2.ActiveAuthGate(); slug != "gate-a" || loaded {
		t.Fatalf("ActiveAuthGate (unloaded): got (%q,%v) want (gate-a,false)", slug, loaded)
	}

	// No gate enabled.
	m3 := seedManager(
		map[string]*DiscoveredEntry{"chatbot": {Slug: "chatbot", Permissions: []string{"chat.send"}}},
		map[string]bool{"chatbot": true},
		map[string]*Loaded{"chatbot": {}},
	)
	if slug, _ := m3.ActiveAuthGate(); slug != "" {
		t.Fatalf("ActiveAuthGate with no gate: got %q want empty", slug)
	}
}
