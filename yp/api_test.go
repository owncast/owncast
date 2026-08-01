package yp

// The directory endpoint is public by design: the Owncast Directory fetches it
// anonymously, with no token, so the viewer-auth gate deliberately never
// challenges it. That makes this handler the only place the "directory is
// switched off" decision is enforced, which is what these tests pin down.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
)

// stubConfig implements only the getters the YP paths touch. The embedded
// interface is nil on purpose: any other call panics loudly instead of
// silently returning a zero value.
type stubConfig struct {
	configrepository.ConfigRepository
	directoryEnabled bool
}

func (s stubConfig) GetDirectoryEnabled() bool               { return s.directoryEnabled }
func (s stubConfig) GetServerName() string                   { return "Test Server" }
func (s stubConfig) GetServerSummary() string                { return "summary" }
func (s stubConfig) GetStreamTitle() string                  { return "title" }
func (s stubConfig) GetNSFW() bool                           { return false }
func (s stubConfig) GetServerMetadataTags() []string         { return nil }
func (s stubConfig) GetSocialHandles() []models.SocialHandle { return nil }
func (s stubConfig) GetServerURL() string                    { return "https://example.com" }

func newYP(directoryEnabled bool, directoryAvailable func() bool) *YP {
	y := New(Deps{
		GetStatus:        func() models.Status { return models.Status{Online: true} },
		ConfigRepository: stubConfig{directoryEnabled: directoryEnabled},
	})
	if directoryAvailable != nil {
		y.SetDirectoryAvailable(directoryAvailable)
	}
	return y
}

func getYP(t *testing.T, y *YP) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	y.GetYPResponse(w, httptest.NewRequest(http.MethodGet, "/api/yp", nil))
	return w
}

func TestGetYPResponse_ServedWhenDirectoryEnabledAndAvailable(t *testing.T) {
	w := getYP(t, newYP(true, func() bool { return true }))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body ypDetailsResponse
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Name != "Test Server" {
		t.Fatalf("name = %q, want Test Server", body.Name)
	}
}

// This is the enforcement that moved out of the gate middleware: when the
// operator blocks public stream status, the gate stops challenging /api/yp and
// this handler is what refuses to answer.
func TestGetYPResponse_404WhenTheGateBlocksTheDirectory(t *testing.T) {
	w := getYP(t, newYP(true, func() bool { return false }))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("blocked directory leaked a body: %q", w.Body.String())
	}
}

func TestGetYPResponse_404WhenDirectoryListingIsOff(t *testing.T) {
	w := getYP(t, newYP(false, func() bool { return true }))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

// The plugin host is optional: with no host wired, the policy hook is nil and
// the directory must behave exactly as it did before gates existed.
func TestGetYPResponse_NoPluginHostLeavesDirectoryWorking(t *testing.T) {
	if w := getYP(t, newYP(true, nil)); w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if w := getYP(t, newYP(false, nil)); w.Code != http.StatusNotFound {
		t.Fatalf("directory off: status = %d, want 404", w.Code)
	}
}

// A blocked directory must also stop the outbound ping, otherwise the instance
// keeps announcing itself to a directory that can no longer read it.
func TestPing_SkippedWhenTheGateBlocksTheDirectory(t *testing.T) {
	tests := []struct {
		name               string
		directoryEnabled   bool
		directoryAvailable func() bool
	}{
		{"gate blocks the directory", true, func() bool { return false }},
		{"directory listing is off", false, func() bool { return true }},
		{"both off", false, func() bool { return false }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			y := newYP(tc.directoryEnabled, tc.directoryAvailable)
			// GetServerURL is only reached once both checks pass; the stub
			// would otherwise send a real request. Reaching it means the
			// guard failed to stop the ping.
			y.getStatus = func() models.Status {
				t.Fatal("ping proceeded past the directory-availability guard")
				return models.Status{}
			}
			y.ping()
		})
	}
}
