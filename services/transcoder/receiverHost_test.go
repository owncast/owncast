package transcoder

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/owncast/owncast/config"
)

type noopReceiverCallbacks struct{}

func (noopReceiverCallbacks) SegmentWritten(string)         {}
func (noopReceiverCallbacks) VariantPlaylistWritten(string) {}
func (noopReceiverCallbacks) MasterPlaylistWritten(string)  {}

func assertNumericPort(t *testing.T, port string) {
	t.Helper()
	if port == "" {
		t.Fatal("listener port was not recorded")
	}
	if _, err := strconv.Atoi(port); err != nil {
		t.Errorf("listener port %q is not numeric: %v", port, err)
	}
}

// TestReceiverBindHost covers the configurable receiver bind. An empty host
// must fall back to loopback rather than binding every interface, and the
// chosen port must be parsed correctly even for a bracketed IPv6 listener
// address (the old strings.Split(":") parse returned an empty port there).
func TestReceiverBindHost(t *testing.T) {
	t.Run("empty host defaults to loopback and records the port", func(t *testing.T) {
		cfg := &config.Config{}
		NewFileWriterReceiverService(cfg).SetupFileWriterReceiverService(noopReceiverCallbacks{})
		assertNumericPort(t, cfg.InternalHLSListenerPort)
	})

	t.Run("ipv6 host records a numeric port", func(t *testing.T) {
		probe, err := net.Listen("tcp", "[::1]:0")
		if err != nil {
			t.Skip("no IPv6 loopback available")
		}
		_ = probe.Close()

		cfg := &config.Config{InternalHLSListenerHost: "::1"}
		NewFileWriterReceiverService(cfg).SetupFileWriterReceiverService(noopReceiverCallbacks{})
		assertNumericPort(t, cfg.InternalHLSListenerPort)
	})
}

type receiverCallbacks struct {
	segments         []string
	variantPlaylists []string
	masterPlaylists  []string
}

func (c *receiverCallbacks) SegmentWritten(path string) {
	c.segments = append(c.segments, path)
}

func (c *receiverCallbacks) VariantPlaylistWritten(path string) {
	c.variantPlaylists = append(c.variantPlaylists, path)
}

func (c *receiverCallbacks) MasterPlaylistWritten(path string) {
	c.masterPlaylists = append(c.masterPlaylists, path)
}

func withTestHLSStoragePath(t *testing.T) string {
	t.Helper()

	originalPath := config.HLSStoragePath
	basePath := filepath.Join(t.TempDir(), "hls")
	if err := os.MkdirAll(filepath.Join(basePath, "0"), 0o750); err != nil {
		t.Fatal(err)
	}

	config.HLSStoragePath = basePath
	t.Cleanup(func() {
		config.HLSStoragePath = originalPath
	})

	return basePath
}

func TestReceiverUploadWritesOnlySafeHLSPaths(t *testing.T) {
	basePath := withTestHLSStoragePath(t)
	callbacks := &receiverCallbacks{}
	receiver := &FileWriterReceiverService{callbacks: callbacks}

	cases := []struct {
		target    string
		body      string
		writePath string
	}{
		{"/stream.m3u8", "master", filepath.Join(basePath, "stream.m3u8")},
		{"/0/stream.m3u8", "variant", filepath.Join(basePath, "0", "stream.m3u8")},
		{"/0/segment.ts", "segment", filepath.Join(basePath, "0", "segment.ts")},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodPut, tc.target, strings.NewReader(tc.body))
		res := httptest.NewRecorder()

		receiver.uploadHandler(res, req)

		if res.Code != http.StatusOK {
			t.Fatalf("expected upload status %d for %s, got %d", http.StatusOK, tc.target, res.Code)
		}

		data, err := os.ReadFile(tc.writePath)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != tc.body {
			t.Fatalf("expected uploaded data %q at %s, got %q", tc.body, tc.writePath, string(data))
		}
	}

	if len(callbacks.masterPlaylists) != 1 || callbacks.masterPlaylists[0] != cases[0].writePath {
		t.Fatalf("expected master playlist callback for %q, got %#v", cases[0].writePath, callbacks.masterPlaylists)
	}
	if len(callbacks.variantPlaylists) != 1 || callbacks.variantPlaylists[0] != cases[1].writePath {
		t.Fatalf("expected variant playlist callback for %q, got %#v", cases[1].writePath, callbacks.variantPlaylists)
	}
	if len(callbacks.segments) != 1 || callbacks.segments[0] != cases[2].writePath {
		t.Fatalf("expected segment callback for %q, got %#v", cases[2].writePath, callbacks.segments)
	}
}

func TestSafeHLSUploadPath(t *testing.T) {
	basePath := withTestHLSStoragePath(t)

	validCases := map[string]string{
		"/stream.m3u8":       filepath.Join(basePath, "stream.m3u8"),
		"/0/stream.m3u8":     filepath.Join(basePath, "0", "stream.m3u8"),
		"/0/segment.ts":      filepath.Join(basePath, "0", "segment.ts"),
		"/audio/segment.aac": filepath.Join(basePath, "audio", "segment.aac"),
		"/v0/part%200.ts":    filepath.Join(basePath, "v0", "part 0.ts"),
	}

	for target, want := range validCases {
		t.Run("valid "+target, func(t *testing.T) {
			_, got, err := safeHLSUploadTarget(basePath, target)
			if err != nil {
				t.Fatalf("expected valid path, got error: %v", err)
			}
			if got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}
		})
	}

	invalidCases := []string{
		"",
		"0/segment.ts",
		"/",
		"//tmp/evil.ts",
		"/../evil.ts",
		"/%2e%2e/evil.ts",
		"/%2E%2E/evil.ts",
		"/..%2fevil.ts",
		"/0/../evil.ts",
		"/0/%2e%2e/evil.ts",
		"/0/%2e%2e%2fevil.ts",
		"/0//segment.ts",
		"/0/./segment.ts",
		"/0/%2e/segment.ts",
		"/0%5c..%5cevil.ts",
		"/0\\..\\evil.ts",
		"/%2Ftmp%2Fevil.ts",
		"/bad%ZZpath.ts",
	}

	for _, target := range invalidCases {
		t.Run("invalid "+target, func(t *testing.T) {
			if relativePath, got, err := safeHLSUploadTarget(basePath, target); err == nil {
				t.Fatalf("expected error for %q, got relative path %q and path %q", target, relativePath, got)
			}
		})
	}
}

func TestSafeHLSUploadPathKeepsOriginalBase(t *testing.T) {
	_, got, err := safeHLSUploadTarget(filepath.Join("data", "hls"), "/0/segment.ts")
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join("data", "hls", "0", "segment.ts")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestReceiverUploadCannotEscapeThroughSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires privileges on Windows")
	}

	basePath := withTestHLSStoragePath(t)
	outsidePath := filepath.Join(t.TempDir(), "outside")
	if err := os.MkdirAll(outsidePath, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsidePath, filepath.Join(basePath, "link")); err != nil {
		t.Skipf("unable to create symlink: %v", err)
	}

	callbacks := &receiverCallbacks{}
	receiver := &FileWriterReceiverService{callbacks: callbacks}
	req := httptest.NewRequest(http.MethodPut, "/link/evil.ts", strings.NewReader("segment"))
	res := httptest.NewRecorder()

	receiver.uploadHandler(res, req)

	if res.Code == http.StatusOK {
		t.Fatalf("expected symlink escape upload to fail")
	}
	if _, err := os.Stat(filepath.Join(outsidePath, "evil.ts")); !os.IsNotExist(err) {
		t.Fatalf("unsafe upload created outside file: %v", err)
	}
	if len(callbacks.masterPlaylists) != 0 || len(callbacks.variantPlaylists) != 0 || len(callbacks.segments) != 0 {
		t.Fatalf("unsafe upload fired callbacks: %#v", callbacks)
	}
}

func TestReceiverUploadRejectsUnsafePaths(t *testing.T) {
	basePath := withTestHLSStoragePath(t)
	callbacks := &receiverCallbacks{}
	receiver := &FileWriterReceiverService{callbacks: callbacks}

	cases := []string{
		"/../evil.ts",
		"/%2e%2e/evil.ts",
		"/0/../evil.ts",
		"/0/%2e%2e/evil.ts",
		"/%2Ftmp%2Fevil.ts",
		"/0//segment.ts",
		"/0/./segment.ts",
		"/0%5c..%5cevil.ts",
	}

	for _, target := range cases {
		t.Run(target, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, target, strings.NewReader("segment"))
			res := httptest.NewRecorder()

			receiver.uploadHandler(res, req)

			if res.Code != http.StatusBadRequest {
				t.Fatalf("expected upload status %d, got %d", http.StatusBadRequest, res.Code)
			}
			if _, err := os.Stat(filepath.Join(filepath.Dir(basePath), "evil.ts")); !os.IsNotExist(err) {
				t.Fatalf("unsafe upload created outside file: %v", err)
			}
			if len(callbacks.masterPlaylists) != 0 || len(callbacks.variantPlaylists) != 0 || len(callbacks.segments) != 0 {
				t.Fatalf("unsafe upload fired callbacks: %#v", callbacks)
			}
		})
	}
}

func TestReceiverUploadRejectsNonPutMethods(t *testing.T) {
	withTestHLSStoragePath(t)
	callbacks := &receiverCallbacks{}
	receiver := &FileWriterReceiverService{callbacks: callbacks}

	req := httptest.NewRequest(http.MethodPost, "/0/segment.ts", strings.NewReader("segment"))
	res := httptest.NewRecorder()

	receiver.uploadHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected upload status %d, got %d", http.StatusBadRequest, res.Code)
	}
	if len(callbacks.masterPlaylists) != 0 || len(callbacks.variantPlaylists) != 0 || len(callbacks.segments) != 0 {
		t.Fatalf("rejected upload fired callbacks: %#v", callbacks)
	}
}
