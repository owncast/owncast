package plugins

import (
	"errors"
	"net/http"

	extism "github.com/extism/go-sdk"
	"github.com/gobwas/glob"
)

// extism's built-in http_request sends plugin fetches through
// http.DefaultClient and checks network.allowedHosts only against the initial
// URL. Re-check every redirect hop against the calling plugin's allowlist so a
// permitted host cannot bounce the request somewhere the plugin isn't allowed
// to reach. Requests without a plugin in their context keep Go's default
// redirect policy. Don't replace http.DefaultClient anywhere in the process;
// doing so drops this check.
func init() {
	http.DefaultClient.CheckRedirect = checkPluginRedirect
}

func checkPluginRedirect(req *http.Request, via []*http.Request) error {
	p, ok := req.Context().Value(extism.PluginCtxKey("plugin")).(*extism.Plugin)
	if !ok {
		// Go's default policy when CheckRedirect is nil.
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	if len(via) >= 10 || !hostAllowed(p.AllowedHosts, req.URL.Hostname()) {
		// Hand the plugin the 3xx unfollowed instead of failing the call: an
		// error here panics inside the extism host function, which leaks the
		// guest shadow stack and forces a plugin reload.
		return http.ErrUseLastResponse
	}
	return nil
}

// hostAllowed mirrors extism's initial allowedHosts check: exact match or
// gobwas glob match against the URL hostname.
func hostAllowed(allowed []string, host string) bool {
	for _, pattern := range allowed {
		if pattern == host {
			return true
		}
		if g, err := glob.Compile(pattern); err == nil && g.Match(host) {
			return true
		}
	}
	return false
}
