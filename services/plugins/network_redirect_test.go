package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	extism "github.com/extism/go-sdk"
)

// A plugin allowed to reach one host must not be able to follow a redirect
// from that host to one outside its allowlist, while redirects that stay
// inside the allowlist keep working. Each script fetches the outside and inside
// redirects, then POSTs the results to the collect URL.
func TestNetworkFetchRedirectRespectsAllowedHosts(t *testing.T) {
	scripts := map[string]func(outside, inside, collect string) string{
		RuntimeJavaScript: func(outside, inside, collect string) string {
			return `const { definePlugin, owncast } = require("@owncast/plugin-sdk");
module.exports = definePlugin({
  onChatMessage() {
    const out = owncast.http.fetch(` + strconv.Quote(outside) + `);
    const ins = owncast.http.fetch(` + strconv.Quote(inside) + `);
    owncast.http.fetch(` + strconv.Quote(collect) + `, {
      method: "POST",
      body: JSON.stringify({ outStatus: out.status, outBody: out.body, inStatus: ins.status, inBody: ins.body }),
    });
  },
});
`
		},
		RuntimePython: func(outside, inside, collect string) string {
			return `import json

@plugin.on_chat_message
def probe(msg):
    out = owncast.http.fetch(` + strconv.Quote(outside) + `)
    ins = owncast.http.fetch(` + strconv.Quote(inside) + `)
    owncast.http.fetch(` + strconv.Quote(collect) + `, {
        "method": "POST",
        "body": json.dumps({"outStatus": out.status, "outBody": out.body, "inStatus": ins.status, "inBody": ins.body}),
    })
`
		},
	}
	for runtime, script := range scripts {
		t.Run(runtime, func(t *testing.T) { testRedirectScope(t, runtime, script) })
	}
}

func testRedirectScope(t *testing.T, runtime string, script func(outside, inside, collect string) string) {
	ctx := context.Background()
	compiledEngines.resetForTest(ctx)
	t.Cleanup(func() { compiledEngines.resetForTest(ctx) })

	const private = "not-for-plugins"
	var outsideHits atomic.Int32
	outside := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outsideHits.Add(1)
		_, _ = w.Write([]byte(private))
	}))
	defer outside.Close()

	var mu sync.Mutex
	var collected string
	allowed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/to-outside":
			http.Redirect(w, r, outside.URL, http.StatusFound)
		case "/to-inside":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			_, _ = w.Write([]byte("inside"))
		case "/collect":
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			collected = string(body)
			mu.Unlock()
		default:
			http.NotFound(w, r)
		}
	}))
	defer allowed.Close()

	// The httptest servers both listen on 127.0.0.1; addressing the allowed
	// one as "localhost" gives the two distinct hostnames.
	u, err := url.Parse(allowed.URL)
	if err != nil {
		t.Fatal(err)
	}
	base := "http://localhost:" + u.Port()

	code := script(base+"/to-outside", base+"/to-inside", base+"/collect")
	manifestBytes, err := json.Marshal(map[string]any{
		"api":         "1",
		"name":        "redirect-scope",
		"slug":        "redirect-scope",
		"version":     "0.1.0",
		"permissions": []string{PermNetworkFetch},
		"network":     map[string]any{"allowedHosts": []string{"localhost"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := loadFromBytes(ctx, &HostEnv{}, manifestBytes, []byte(code), runtime, "redirect-scope", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer loaded.Close(ctx)

	NewDispatcher([]*Loaded{loaded}).Dispatch(ctx, EventChatMessageReceived, chatPayload("alice", "trigger"))

	mu.Lock()
	defer mu.Unlock()
	var got struct {
		OutStatus int    `json:"outStatus"`
		OutBody   string `json:"outBody"`
		InStatus  int    `json:"inStatus"`
		InBody    string `json:"inBody"`
	}
	if err := json.Unmarshal([]byte(collected), &got); err != nil {
		t.Fatalf("plugin did not report results (%q): %v", collected, err)
	}
	if n := outsideHits.Load(); n != 0 || got.OutStatus != http.StatusFound || got.OutBody == private {
		t.Fatalf("redirect outside allowedHosts was followed: hits=%d status=%d body=%q", n, got.OutStatus, got.OutBody)
	}
	if got.InStatus != http.StatusOK || got.InBody != "inside" {
		t.Fatalf("redirect inside allowedHosts was not followed: status=%d body=%q", got.InStatus, got.InBody)
	}
}

func TestCheckPluginRedirectWithoutPlugin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://anywhere.example/", nil)
	if err := checkPluginRedirect(req, nil); err != nil {
		t.Fatalf("non-plugin redirect blocked: %v", err)
	}
	if err := checkPluginRedirect(req, make([]*http.Request, 10)); err == nil {
		t.Fatal("expected the default 10-redirect limit")
	}
}

// Plugin requests must never get a plain error from CheckRedirect (it panics
// inside extism), including at the redirect limit.
func TestCheckPluginRedirectLimitForPlugin(t *testing.T) {
	p := &extism.Plugin{AllowedHosts: []string{"allowed.example"}}
	ctx := context.WithValue(context.Background(), extism.PluginCtxKey("plugin"), p)
	for _, host := range []string{"allowed.example", "other.example"} {
		req := httptest.NewRequest(http.MethodGet, "http://"+host+"/", nil).WithContext(ctx)
		if err := checkPluginRedirect(req, make([]*http.Request, 10)); !errors.Is(err, http.ErrUseLastResponse) {
			t.Errorf("%s at redirect limit: got %v, want http.ErrUseLastResponse", host, err)
		}
	}
}
