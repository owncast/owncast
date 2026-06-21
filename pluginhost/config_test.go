package pluginhost

import (
	"testing"

	"github.com/owncast/owncast/services/plugins/kv"
)

func TestCoerceConfigValue(t *testing.T) {
	cases := []struct {
		name      string
		fieldType string
		in        any
		wantErr   bool
		want      any
	}{
		{"string ok", "string", "hi", false, "hi"},
		{"string from number rejected", "string", float64(5), true, nil},
		{"number ok", "number", float64(3), false, float64(3)},
		{"number from string rejected", "number", "3", true, nil},
		{"boolean ok", "boolean", true, false, true},
		{"boolean from string rejected", "boolean", "true", true, nil},
		{"unknown type passes through", "", map[string]any{"x": 1}, false, map[string]any{"x": 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := coerceConfigValue(tc.fieldType, tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got value %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.fieldType != "" && got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestReadConfigOverrides(t *testing.T) {
	store := kv.NewMemory()

	// No overrides saved yet → nil.
	if got := readConfigOverrides(store, "plug"); got != nil {
		t.Fatalf("expected nil for unset overrides, got %v", got)
	}
	// nil store → nil (no panic).
	if got := readConfigOverrides(nil, "plug"); got != nil {
		t.Fatalf("expected nil for nil store, got %v", got)
	}

	if err := store.Namespace("plug").Set(runtimeConfigKey, []byte(`{"password":"s3cret","count":5,"on":true}`)); err != nil {
		t.Fatal(err)
	}
	m := readConfigOverrides(store, "plug")
	if m["password"] != "s3cret" {
		t.Fatalf("password: got %v", m["password"])
	}
	if m["count"] != float64(5) {
		t.Fatalf("count: got %v (%T)", m["count"], m["count"])
	}
	if m["on"] != true {
		t.Fatalf("on: got %v", m["on"])
	}
	// A different plugin's namespace is isolated.
	if got := readConfigOverrides(store, "other"); got != nil {
		t.Fatalf("expected nil for a different plugin, got %v", got)
	}
}
