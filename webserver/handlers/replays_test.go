package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/replays"
)

func newReplayTestHandlers(t *testing.T) (*Handlers, *replays.Service, configrepository.ConfigRepository) {
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

	return &Handlers{replays: svc, configRepository: cr}, svc, cr
}

// enableReplays turns the feature on through the admin setting, which is how
// operators enable it.
func enableReplays(t *testing.T, cr configrepository.ConfigRepository) {
	t.Helper()

	if err := cr.SetReplayFeaturesEnabled(true); err != nil {
		t.Fatalf("unable to enable replay features: %v", err)
	}
}

// TestReplayEndpointsDisabled verifies the replay endpoints 404 when the
// feature is off, regardless of any recorded data.
func TestReplayEndpointsDisabled(t *testing.T) {
	h, _, _ := newReplayTestHandlers(t)

	cases := map[string]http.HandlerFunc{
		"/api/clips":  h.GetAllClips,
		"/replay/abc": h.GetReplay,
		"/clip/abc":   h.GetClip,
		"/clips/abc":  h.GetClipPage,
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

// TestResolveClipWindow covers the clip-window arithmetic: a trailing window
// resolves against the newest recorded media, and explicit windows are
// validated and clamped.
func TestResolveClipWindow(t *testing.T) {
	const mediaEnd = 100.0
	const maxDuration = 60.0

	t.Run("trailing window ends at the newest media", func(t *testing.T) {
		start, end, err := resolveClipWindow(clipWindowRequest{DurationSeconds: 30}, mediaEnd, maxDuration)
		if err != nil {
			t.Fatal(err)
		}
		if start != 70 || end != 100 {
			t.Errorf("expected window 70-100, got %v-%v", start, end)
		}
	})

	t.Run("trailing window is capped at the operator maximum", func(t *testing.T) {
		start, end, err := resolveClipWindow(clipWindowRequest{DurationSeconds: 600}, mediaEnd, maxDuration)
		if err != nil {
			t.Fatal(err)
		}
		if end-start != maxDuration {
			t.Errorf("expected a %v second window, got %v", maxDuration, end-start)
		}
	})

	t.Run("trailing window never starts before the stream", func(t *testing.T) {
		start, _, err := resolveClipWindow(clipWindowRequest{DurationSeconds: 50}, 20, maxDuration)
		if err != nil {
			t.Fatal(err)
		}
		if start != 0 {
			t.Errorf("expected the window to start at 0, got %v", start)
		}
	})

	t.Run("explicit window is kept verbatim", func(t *testing.T) {
		start, end, err := resolveClipWindow(clipWindowRequest{StartSeconds: 10.5, EndSeconds: 25.25}, mediaEnd, maxDuration)
		if err != nil {
			t.Fatal(err)
		}
		if start != 10.5 || end != 25.25 {
			t.Errorf("expected window 10.5-25.25, got %v-%v", start, end)
		}
	})

	t.Run("explicit window is clamped to recorded media", func(t *testing.T) {
		_, end, err := resolveClipWindow(clipWindowRequest{StartSeconds: 90, EndSeconds: 130}, mediaEnd, maxDuration)
		if err != nil {
			t.Fatal(err)
		}
		if end != mediaEnd {
			t.Errorf("expected the window to end at %v, got %v", mediaEnd, end)
		}
	})

	t.Run("rejects invalid and oversized windows", func(t *testing.T) {
		cases := map[string]clipWindowRequest{
			"negative start":     {StartSeconds: -1, EndSeconds: 10},
			"end before start":   {StartSeconds: 20, EndSeconds: 10},
			"start past the end": {StartSeconds: 150, EndSeconds: 160},
			"longer than max":    {StartSeconds: 0, EndSeconds: 90},
		}

		for name, request := range cases {
			if _, _, err := resolveClipWindow(request, mediaEnd, maxDuration); err == nil {
				t.Errorf("%s: expected an error", name)
			}
		}
	})
}

// TestGetReplayMasterPlaylist verifies the master playlist endpoint returns a
// playable HLS manifest for a recorded stream.
func TestGetReplayMasterPlaylist(t *testing.T) {
	h, svc, cr := newReplayTestHandlers(t)
	enableReplays(t, cr)

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

	// An in-progress recording must not be cached, or players stall on a
	// stale segment list.
	if cc := rec.Header().Get("Cache-Control"); !strings.Contains(cc, "no-cache") {
		t.Errorf("expected an uncacheable in-progress playlist, got %q", cc)
	}

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected CORS to be enabled for players, got %q", origin)
	}

	if !strings.Contains(rec.Body.String(), "#EXTM3U") {
		t.Errorf("expected an HLS playlist, got:\n%s", rec.Body.String())
	}
}

// TestAddClipRejectsUnknownStream verifies clip creation validates its target.
func TestAddClipRejectsUnknownStream(t *testing.T) {
	h, _, cr := newReplayTestHandlers(t)
	enableReplays(t, cr)

	req := httptest.NewRequest(http.MethodPost, "/api/clip", strings.NewReader(`{"streamId":"nope","durationSeconds":30}`))
	rec := httptest.NewRecorder()
	h.AddClip(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for an unknown stream, got %d", rec.Code)
	}
}

// TestAddClipRejectsStreamWithNoVideo verifies a stream that has recorded
// nothing yet cannot be clipped: there is no media timeline to clip against.
func TestAddClipRejectsStreamWithNoVideo(t *testing.T) {
	h, svc, cr := newReplayTestHandlers(t)
	enableReplays(t, cr)

	if recording := svc.NewRecording("teststream3"); recording == nil {
		t.Fatal("expected a recording to be created")
	}

	req := httptest.NewRequest(http.MethodPost, "/api/clip", strings.NewReader(`{"streamId":"teststream3","durationSeconds":30}`))
	rec := httptest.NewRecorder()
	h.AddClip(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for a stream with no recorded video, got %d", rec.Code)
	}
}

// TestAddClipDisabledByOperator verifies the clips-enabled setting is honored
// even while replay features are on.
func TestAddClipDisabledByOperator(t *testing.T) {
	h, _, cr := newReplayTestHandlers(t)
	enableReplays(t, cr)

	if err := cr.SetClipsEnabled(false); err != nil {
		t.Fatalf("unable to disable clips: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/clip", strings.NewReader(`{"streamId":"teststream1","durationSeconds":30}`))
	rec := httptest.NewRecorder()
	h.AddClip(rec, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unable to decode response: %v", err)
	}

	if response["success"] != false {
		t.Errorf("expected clip creation to be refused, got %v", rec.Body.String())
	}
}

// TestClipRateLimit verifies repeated clip requests from one client are
// throttled, since clip creation is open to any viewer.
func TestClipRateLimit(t *testing.T) {
	ip := "203.0.113.77"

	clipRateLimitersMu.Lock()
	delete(clipRateLimiters, ip)
	clipRateLimitersMu.Unlock()

	for i := range clipBurst {
		if !allowClipFromIP(ip) {
			t.Fatalf("expected the first %d clips to be allowed, request %d was denied", clipBurst, i+1)
		}
	}

	if allowClipFromIP(ip) {
		t.Error("expected a clip request beyond the burst to be denied")
	}
}
