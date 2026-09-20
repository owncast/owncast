package admin

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/persistence/migrations"
	"github.com/owncast/owncast/services/datastore"
)

func TestFederatedServerPublicAndAdminPayloads(t *testing.T) {
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := migrations.Run(ds.DB, t.TempDir()); err != nil {
		t.Fatal(err)
	}
	repo := federatedserversrepository.New(ds)
	federatedserversrepository.SetGlobalInstance(repo)

	for _, server := range []struct {
		iri, name string
	}{
		{"https://priority-one.example.com", "Priority One"},
		{"https://priority-two.example.com", "Priority Two"},
	} {
		if err := repo.AddFederatedServer(server.iri, server.name, "https://example.com/logo.png", time.Now(), false, server.name, "accepted"); err != nil {
			t.Fatal(err)
		}
		if err := repo.AssignNextPriority(server.iri); err != nil {
			t.Fatal(err)
		}
	}
	if err := repo.AddFederatedServer("https://pending.example.com", "Pending", "", time.Now(), true, "pending", "pending"); err != nil {
		t.Fatal(err)
	}

	publicRecorder := httptest.NewRecorder()
	(&Admin{}).GetFederatedServers(publicRecorder, httptest.NewRequest("GET", "/api/federation/servers", nil))
	var publicResponse struct {
		Servers []map[string]json.RawMessage `json:"servers"`
	}
	if err := json.Unmarshal(publicRecorder.Body.Bytes(), &publicResponse); err != nil {
		t.Fatal(err)
	}
	if len(publicResponse.Servers) != 2 {
		t.Fatalf("public server count = %d, want 2", len(publicResponse.Servers))
	}
	if string(publicResponse.Servers[0]["name"]) != `"Priority One"` {
		t.Fatalf("public order = %s, want Priority One first", publicResponse.Servers[0]["name"])
	}
	for _, field := range []string{"lastStatusUpdate", "addedAt", "pending", "followStatus", "priority"} {
		if _, ok := publicResponse.Servers[0][field]; ok {
			t.Errorf("public payload unexpectedly contains %q", field)
		}
	}
	if _, ok := publicResponse.Servers[0]["iri"]; !ok {
		t.Error("public payload omitted viewer field iri")
	}

	adminRecorder := httptest.NewRecorder()
	(&Admin{}).GetAdminFederatedServers(adminRecorder, httptest.NewRequest("GET", "/api/admin/federation/servers", nil))
	var adminResponse struct {
		Servers []map[string]json.RawMessage `json:"servers"`
	}
	if err := json.Unmarshal(adminRecorder.Body.Bytes(), &adminResponse); err != nil {
		t.Fatal(err)
	}
	if len(adminResponse.Servers) != 3 {
		t.Fatalf("admin server count = %d, want 3", len(adminResponse.Servers))
	}
	for _, field := range []string{"addedAt", "pending", "followStatus", "priority"} {
		if _, ok := adminResponse.Servers[0][field]; !ok {
			t.Errorf("admin payload omitted management field %q", field)
		}
	}
}
