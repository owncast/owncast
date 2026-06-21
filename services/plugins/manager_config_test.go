package plugins

import "testing"

func TestManager_ConfigSchema(t *testing.T) {
	m := seedManager(
		map[string]*DiscoveredEntry{
			"with-config": {
				Slug: "with-config",
				Config: map[string]ConfigField{
					"password": {Type: "string", Default: "letmein", Description: "Shared password"},
					"limit":    {Type: "number", Default: float64(10)},
				},
			},
			"no-config": {Slug: "no-config"},
		},
		map[string]bool{},
		nil,
	)

	sc := m.ConfigSchema("with-config")
	if sc == nil {
		t.Fatal("expected a config schema")
	}
	if sc["password"].Type != "string" || sc["password"].Default != "letmein" {
		t.Fatalf("password field wrong: %+v", sc["password"])
	}
	if sc["limit"].Type != "number" {
		t.Fatalf("limit field wrong: %+v", sc["limit"])
	}

	if m.ConfigSchema("no-config") != nil {
		t.Fatal("plugin with no config should return nil schema")
	}
	if m.ConfigSchema("unknown") != nil {
		t.Fatal("unknown plugin should return nil schema")
	}
}
