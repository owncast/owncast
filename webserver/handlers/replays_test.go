package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

// TestResolveClipWindow covers explicit media-time window validation.
func TestResolveClipWindow(t *testing.T) {
	const mediaEnd = 100.0
	const maxDuration = 60.0

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

	t.Run("rejects invalid windows", func(t *testing.T) {
		cases := map[string]clipWindowRequest{
			"negative start":   {StartSeconds: -1, EndSeconds: 10},
			"end before start": {StartSeconds: 20, EndSeconds: 10},
			"start past the end": {
				StartSeconds: 150,
				EndSeconds:   160,
			},
		}

		for name, request := range cases {
			if _, _, err := resolveClipWindow(request, mediaEnd, maxDuration); err == nil {
				t.Errorf("%s: expected an error", name)
			}
		}
	})

	t.Run("over-long window keeps its end", func(t *testing.T) {
		start, end, err := resolveClipWindow(clipWindowRequest{StartSeconds: 0, EndSeconds: 90}, mediaEnd, maxDuration)
		if err != nil {
			t.Fatal(err)
		}
		if start != 30 || end != 90 {
			t.Errorf("expected window 30-90, got %v-%v", start, end)
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

	req := httptest.NewRequest(http.MethodPost, "/api/clip", strings.NewReader(`{"streamId":"nope","relativeStartTimeSeconds":0,"relativeEndTimeSeconds":30}`))
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

	req := httptest.NewRequest(http.MethodPost, "/api/clip", strings.NewReader(`{"streamId":"teststream3","relativeStartTimeSeconds":0,"relativeEndTimeSeconds":30}`))
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

	req := httptest.NewRequest(http.MethodPost, "/api/clip", strings.NewReader(`{"streamId":"teststream1","relativeStartTimeSeconds":0,"relativeEndTimeSeconds":30}`))
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

// TestClipPermissionDenialReason covers who may create clips at each
// permission level: identity required everywhere, moderators always allowed,
// authentication and account age gates enforced.
func TestClipPermissionDenialReason(t *testing.T) {
	oldUser := &models.User{CreatedAt: time.Now().Add(-2 * time.Hour), DisplayName: "old"}
	newUser := &models.User{CreatedAt: time.Now().Add(-5 * time.Minute), DisplayName: "new"}
	authedNewUser := &models.User{CreatedAt: time.Now().Add(-5 * time.Minute), Authenticated: true}
	moderator := &models.User{CreatedAt: time.Now(), Scopes: []string{"MODERATOR"}}
	disabledAt := time.Now()
	banned := &models.User{CreatedAt: time.Now().Add(-2 * time.Hour), DisabledAt: &disabledAt}

	cases := []struct {
		name        string
		user        *models.User
		permissions string
		allowed     bool
	}{
		{"anonymous is always rejected", nil, models.ClipPermissionsEstablished, false},
		{"banned is always rejected", banned, models.ClipPermissionsEstablished, false},
		{"established allows old accounts", oldUser, models.ClipPermissionsEstablished, true},
		{"established rejects young accounts", newUser, models.ClipPermissionsEstablished, false},
		{"authenticated rejects unauthenticated", oldUser, models.ClipPermissionsAuthenticated, false},
		{"authenticated allows authenticated", authedNewUser, models.ClipPermissionsAuthenticated, true},
		{"moderators rejects non-moderators", oldUser, models.ClipPermissionsModerators, false},
		{"moderators allows moderators", moderator, models.ClipPermissionsModerators, true},
		{"moderator bypasses account age", moderator, models.ClipPermissionsEstablished, true},
		{"moderator bypasses authentication", moderator, models.ClipPermissionsAuthenticated, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reason := clipPermissionDenialReason(tc.user, tc.permissions)
			if tc.allowed && reason != "" {
				t.Errorf("expected allowed, got %q", reason)
			}
			if !tc.allowed && reason == "" {
				t.Error("expected a denial reason")
			}
		})
	}
}
