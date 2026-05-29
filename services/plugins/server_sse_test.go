package plugins

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gobwas/glob"
)

// sseServer builds a Server with a single "demo" plugin and the given hub.
// adminGlobs (if any) are compiled into the plugin's admin-path set so the
// reserved _sse endpoint can be exercised under admin gating.
func sseServer(t *testing.T, perms []string, hub *SSEHub, adminGlobs ...string) *Server {
	t.Helper()
	var globs []glob.Glob
	for _, g := range adminGlobs {
		globs = append(globs, glob.MustCompile(g))
	}
	loaded := &Loaded{
		Manifest:   &Manifest{API: "1", DisplayName: "demo", Slug: "demo", Version: "1.0.0", Permissions: perms},
		adminGlobs: globs,
	}
	s := NewServer([]*Loaded{loaded})
	s.SSE = hub
	return s
}

func TestServer_SSE_NotFoundWithoutPermission(t *testing.T) {
	// http.serve alone does not grant the reserved SSE endpoint.
	s := sseServer(t, []string{"http.serve"}, NewSSEHub())
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/_sse/overlay", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("without http.sse: status %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestServer_SSE_ServiceUnavailableWhenHubNil(t *testing.T) {
	// Permission granted but the host wired no hub → 503, not a panic.
	s := sseServer(t, []string{"http.sse"}, nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/_sse/overlay", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("nil hub: status %d want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestServer_SSE_ServiceUnavailableAtConnectionCap(t *testing.T) {
	hub := NewSSEHub()
	hub.maxPerPlugin = 0 // every connection is over the cap
	s := sseServer(t, []string{"http.sse"}, hub)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/_sse/overlay", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("over cap: status %d want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestServer_SSE_AdminChannelRequiresAuth(t *testing.T) {
	// A channel matching an admin glob is auth-gated like any admin path.
	// The server delegates the 401 to the host's RequireAdmin middleware
	// (the real one in production is middleware.RequireAdminAuth); here
	// we stub it with a always-reject wrapper so we observe the 401.
	s := sseServer(t, []string{"http.sse"}, NewSSEHub(), "/_sse/*")
	s.RequireAdmin = func(http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/_sse/admin-stats", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("admin SSE channel without auth: status %d want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestServer_SSE_StreamsPublishedFrames(t *testing.T) {
	// Both the named-channel and default-channel routing forms should connect
	// and stream a frame the plugin publishes to the matching channel.
	cases := []struct{ name, path, channel string }{
		{"named channel", "/_sse/overlay", "overlay"},
		{"default channel", "/_sse", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hub := NewSSEHub()
			srv := httptest.NewServer(sseServer(t, []string{"http.sse"}, hub))
			defer srv.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "GET", srv.URL+"/plugins/demo"+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("connect: status %d want %d", resp.StatusCode, http.StatusOK)
			}
			if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
				t.Errorf("Content-Type = %q want text/event-stream", ct)
			}

			// The handler subscribes on its own goroutine; wait for it before
			// publishing so the frame isn't dropped before anyone's listening.
			waitForSSESubscribers(t, hub, "demo", tc.channel, 1)

			if n := hub.Publish("demo", tc.channel, "emoji", []byte("🦉")); n != 1 {
				t.Fatalf("expected delivery to the connected client, got %d", n)
			}

			frame := readSSEFrame(t, resp.Body)
			if !strings.Contains(frame, "event: emoji\n") || !strings.Contains(frame, "data: 🦉\n") {
				t.Errorf("unexpected SSE frame: %q", frame)
			}

			// Client disconnect must release the subscription.
			cancel()
			waitForSSESubscribers(t, hub, "demo", tc.channel, 0)
		})
	}
}

// waitForSSESubscribers polls the hub until (plugin, channel) has exactly want
// subscribers, failing the test on timeout.
func waitForSSESubscribers(t *testing.T, h *SSEHub, plugin, channel string, want int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		h.mu.Lock()
		got := len(h.subscribers[sseKey(plugin, channel)])
		h.mu.Unlock()
		if got == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %d subscriber(s) on %s/%s (have %d)", want, plugin, channel, got)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// readSSEFrame reads from an event-stream body until one complete data frame
// (terminated by a blank line) is read, skipping any keep-alive comments.
func readSSEFrame(t *testing.T, body io.Reader) string {
	t.Helper()
	type result struct {
		frame string
		err   error
	}
	done := make(chan result, 1)
	go func() {
		r := bufio.NewReader(body)
		var b strings.Builder
		for {
			line, err := r.ReadString('\n')
			b.WriteString(line)
			if err != nil {
				done <- result{b.String(), err}
				return
			}
			if line == "\n" && strings.Contains(b.String(), "data:") {
				done <- result{b.String(), nil}
				return
			}
		}
	}()
	select {
	case res := <-done:
		if res.err != nil && !strings.Contains(res.frame, "data:") {
			t.Fatalf("reading SSE frame: %v (got %q)", res.err, res.frame)
		}
		return res.frame
	case <-time.After(3 * time.Second):
		t.Fatal("timed out reading SSE frame")
		return ""
	}
}
