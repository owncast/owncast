package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/replays"
)

// newReplayAdmin builds an *Admin wired to the replay service, with replay
// features enabled through the admin setting.
func newReplayAdmin(t *testing.T) (*Admin, *replays.Service, configrepository.ConfigRepository) {
	t.Helper()

	cr := configrepository.New(testDatastore)

	// A concrete (non-passthrough) variant so recordings produce valid
	// playlists.
	if err := cr.SetStreamOutputVariants([]models.StreamOutputVariant{
		{Name: "high", VideoBitrate: 1200, Framerate: 30, ScaledWidth: 1280, ScaledHeight: 720},
	}); err != nil {
		t.Fatalf("unable to set output variants: %v", err)
	}

	if err := cr.SetReplayFeaturesEnabled(true); err != nil {
		t.Fatalf("unable to enable replay features: %v", err)
	}

	svc := replays.New(testDatastore, cr)

	return &Admin{configRepository: cr, replays: svc}, svc, cr
}

// TestAdminReplayEndpointsDisabled verifies the admin replay surface 404s when
// the feature is off, so the admin UI can hide the section.
func TestAdminReplayEndpointsDisabled(t *testing.T) {
	cr := configrepository.New(testDatastore)
	if err := cr.SetReplayFeaturesEnabled(false); err != nil {
		t.Fatalf("unable to disable replay features: %v", err)
	}
	a := &Admin{configRepository: cr, replays: replays.New(testDatastore, cr)}

	cases := map[string]http.HandlerFunc{
		"/api/admin/replays": a.GetAdminReplays,
		"/api/admin/clips":   a.GetAdminClips,
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

// TestAdminReplayListingAndDelete verifies the admin can list replays with
// their clip counts and delete a clip without touching its replay, then delete
// the replay itself.
func TestAdminReplayListingAndDelete(t *testing.T) {
	a, svc, _ := newReplayAdmin(t)

	streamID := "admin-replay-1"
	recording := svc.NewRecording(streamID)
	if recording == nil {
		t.Fatal("expected a recording to be created")
	}

	// A clip needs recorded media to exist, so report a segment and give it
	// real timing the way the transcoder does.
	writeTestSegments(t, svc, recording, streamID)

	clipID, _, err := svc.AddClipForStream(streamID, "admin clip", "", 0, 2)
	if err != nil {
		t.Fatalf("unable to create clip: %v", err)
	}

	// The replay listing should report the clip against the replay.
	req := httptest.NewRequest(http.MethodGet, "/api/admin/replays", nil)
	rec := httptest.NewRecorder()
	a.GetAdminReplays(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var replayList []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &replayList); err != nil {
		t.Fatalf("unable to decode replays: %v", err)
	}

	found := false
	for _, r := range replayList {
		if r["id"] == streamID {
			found = true
			if count, ok := r["clipCount"].(float64); !ok || count != 1 {
				t.Errorf("expected clipCount 1, got %v", r["clipCount"])
			}
			if bytes, ok := r["totalBytes"].(float64); !ok || bytes <= 0 {
				t.Errorf("expected the replay to report disk usage, got %v", r["totalBytes"])
			}
		}
	}
	if !found {
		t.Fatalf("expected replay %q to be listed", streamID)
	}

	// Deleting the clip leaves the replay in place.
	deleteRec := httptest.NewRecorder()
	a.DeleteClip(deleteRec, httptest.NewRequest(http.MethodPost, "/api/admin/clips/delete", strings.NewReader(`{"id":"`+clipID+`"}`)))
	assertSuccess(t, deleteRec, true)

	clips, err := svc.GetAllClipsForStream(streamID)
	if err != nil {
		t.Fatalf("unable to list clips: %v", err)
	}
	if len(clips) != 0 {
		t.Errorf("expected the clip to be deleted, %d remain", len(clips))
	}

	if _, err := svc.NewPlaylistGenerator().GetStream(streamID); err != nil {
		t.Errorf("expected the replay to survive deleting its clip: %v", err)
	}

	// Deleting the replay removes it entirely.
	replayRec := httptest.NewRecorder()
	a.DeleteReplay(replayRec, httptest.NewRequest(http.MethodPost, "/api/admin/replays/delete", strings.NewReader(`{"id":"`+streamID+`"}`)))
	assertSuccess(t, replayRec, true)

	if _, err := svc.NewPlaylistGenerator().GetStream(streamID); err == nil {
		t.Error("expected the replay to be deleted")
	}
}

// TestAdminDeleteReplayCascadesToClips verifies deleting a replay also deletes
// the clips taken from it, since they reference its recorded video.
func TestAdminDeleteReplayCascadesToClips(t *testing.T) {
	a, svc, _ := newReplayAdmin(t)

	streamID := "admin-replay-2"
	recording := svc.NewRecording(streamID)
	if recording == nil {
		t.Fatal("expected a recording to be created")
	}
	writeTestSegments(t, svc, recording, streamID)

	if _, _, err := svc.AddClipForStream(streamID, "doomed clip", "", 0, 2); err != nil {
		t.Fatalf("unable to create clip: %v", err)
	}

	rec := httptest.NewRecorder()
	a.DeleteReplay(rec, httptest.NewRequest(http.MethodPost, "/api/admin/replays/delete", strings.NewReader(`{"id":"`+streamID+`"}`)))
	assertSuccess(t, rec, true)

	clips, err := svc.GetAllClipsForStream(streamID)
	if err != nil {
		t.Fatalf("unable to list clips: %v", err)
	}
	if len(clips) != 0 {
		t.Errorf("expected clips to be deleted with their replay, %d remain", len(clips))
	}
}

// TestSetMaxClipDurationValidation verifies the configured clip limit is bounded.
func TestSetMaxClipDurationValidation(t *testing.T) {
	a, _, cr := newReplayAdmin(t)

	rec := httptest.NewRecorder()
	a.SetMaxClipDuration(rec, httptest.NewRequest(http.MethodPost, "/api/admin/config/replay/maxclipduration", strings.NewReader(`{"value":0}`)))
	assertSuccess(t, rec, false)

	rec = httptest.NewRecorder()
	a.SetMaxClipDuration(rec, httptest.NewRequest(http.MethodPost, "/api/admin/config/replay/maxclipduration", strings.NewReader(`{"value":99999}`)))
	assertSuccess(t, rec, false)

	rec = httptest.NewRecorder()
	a.SetMaxClipDuration(rec, httptest.NewRequest(http.MethodPost, "/api/admin/config/replay/maxclipduration", strings.NewReader(`{"value":45}`)))
	assertSuccess(t, rec, true)

	if got := cr.GetMaxClipDurationSeconds(); got != 45 {
		t.Errorf("expected max clip duration 45, got %d", got)
	}
}

// writeTestSegments reports a couple of segments for a recording and supplies
// their real durations through a variant playlist, mirroring the transcoder.
func writeTestSegments(t *testing.T, svc *replays.Service, recording *replays.HLSRecorder, streamID string) {
	t.Helper()

	variantDir := filepath.Join(t.TempDir(), "0")
	if err := os.MkdirAll(variantDir, 0o750); err != nil {
		t.Fatalf("unable to create variant directory: %v", err)
	}

	playlist := "#EXTM3U\n#EXT-X-VERSION:6\n#EXT-X-TARGETDURATION:3\n#EXT-X-MEDIA-SEQUENCE:0\n"
	for i := range 2 {
		name := fmt.Sprintf("%s-%d.ts", streamID, i)
		recording.SegmentWritten(filepath.Join("hls", "0", name), 2048)
		playlist += fmt.Sprintf("#EXTINF:2.000,\n%s\n", name)
	}

	playlistPath := filepath.Join(variantDir, "stream.m3u8")
	if err := os.WriteFile(playlistPath, []byte(playlist), 0o600); err != nil {
		t.Fatalf("unable to write variant playlist: %v", err)
	}

	recording.VariantPlaylistWritten(playlistPath)

	if end, err := svc.GetStreamMediaEnd(streamID); err != nil || end <= 0 {
		t.Fatalf("expected recorded media, got end=%v err=%v", end, err)
	}
}

// assertSuccess checks the success field of a standard admin JSON response.
func assertSuccess(t *testing.T, rec *httptest.ResponseRecorder, want bool) {
	t.Helper()

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unable to decode response %q: %v", rec.Body.String(), err)
	}

	if response["success"] != want {
		t.Errorf("expected success=%v, got %q", want, rec.Body.String())
	}
}
