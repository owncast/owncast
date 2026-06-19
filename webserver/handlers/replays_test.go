package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/replays"
)

func newReplayTestHandlers(t *testing.T) (*Handlers, *replays.Service) {
	t.Helper()

	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("unable to set up persistence: %v", err)
	}

	cr := configrepository.New(ds)

	// Configure a concrete (non-passthrough) output variant so recorded
	// streams produce valid playlists (a passthrough variant has no bitrate
	// and would fail playlist validation).
	if err := cr.SetStreamOutputVariants([]models.StreamOutputVariant{
		{Name: "high", VideoBitrate: 1200, Framerate: 30, ScaledWidth: 1280, ScaledHeight: 720},
	}); err != nil {
		t.Fatalf("unable to set output variants: %v", err)
	}

	svc := replays.New(ds, cr)
	svc.Setup()

	return &Handlers{replays: svc}, svc
}

// TestReplayEndpointsDisabled verifies the replay endpoints 404 when the
// feature flag is off, regardless of any recorded data.
func TestReplayEndpointsDisabled(t *testing.T) {
	config.EnableReplayFeatures = false
	h, _ := newReplayTestHandlers(t)

	cases := map[string]http.HandlerFunc{
		"/api/replays": h.GetReplays,
		"/api/clips":   h.GetAllClips,
		"/replay/abc":  h.GetReplay,
		"/clip/abc":    h.GetClip,
	}

	for path, handler := range cases {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s: expected 404 when replays disabled, got %d", path, rec.Code)
		}
	}
}

// TestGetReplaysEnabled verifies that with the feature enabled, a recorded
// stream is listed by /api/replays.
func TestGetReplaysEnabled(t *testing.T) {
	config.EnableReplayFeatures = true
	defer func() { config.EnableReplayFeatures = false }()

	h, svc := newReplayTestHandlers(t)

	// Record a stream so there's something to list.
	recording := svc.NewRecording("teststream1")
	if recording == nil {
		t.Fatal("expected a recording to be created")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/replays", nil)
	rec := httptest.NewRecorder()
	h.GetReplays(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var streams []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &streams); err != nil {
		t.Fatalf("unable to decode replays response: %v", err)
	}

	found := false
	for _, s := range streams {
		if s["id"] == "teststream1" {
			found = true
			if manifest, ok := s["manifest"].(string); !ok || manifest != "/replay/teststream1" {
				t.Errorf("expected manifest /replay/teststream1, got %v", s["manifest"])
			}
		}
	}
	if !found {
		t.Error("expected recorded stream to be listed in replays")
	}
}

// TestGetReplayMasterPlaylist verifies the master playlist endpoint returns a
// playable HLS manifest for a recorded stream.
func TestGetReplayMasterPlaylist(t *testing.T) {
	config.EnableReplayFeatures = true
	defer func() { config.EnableReplayFeatures = false }()

	h, svc := newReplayTestHandlers(t)
	if recording := svc.NewRecording("teststream2"); recording == nil {
		t.Fatal("expected a recording to be created")
	}

	req := httptest.NewRequest(http.MethodGet, "/replay/teststream2", nil)
	rec := httptest.NewRecorder()
	h.GetReplay(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/x-mpegURL" {
		t.Errorf("expected HLS content type, got %q", ct)
	}

	if !strings.Contains(rec.Body.String(), "#EXTM3U") {
		t.Errorf("expected an HLS playlist, got:\n%s", rec.Body.String())
	}
}

// TestAddClipRequiresPost verifies clip creation only accepts POST.
func TestAddClipRequiresPost(t *testing.T) {
	config.EnableReplayFeatures = true
	defer func() { config.EnableReplayFeatures = false }()

	h, _ := newReplayTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/api/clip", nil)
	rec := httptest.NewRecorder()
	h.AddClip(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-POST clip request, got %d", rec.Code)
	}
}
