package handlers

import (
	"net/http"
)

// ServeCustomJavascript serves the admin's custom JavaScript followed
// by the concatenated JS contributed by every loaded plugin's
// manifest.scripts entries (each preceded by a `// plugin: <slug>
// ...` delimiter for devtools attribution). The viewer loads a single
// <script> tag pointing at this endpoint so plugins can extend the
// page without each plugin needing its own <script> tag.
func (h *Handlers) ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	body := []byte(h.configRepository.GetCustomJavascript())
	if h.pluginJSContent != nil {
		pluginJS := h.pluginJSContent()
		if len(pluginJS) > 0 {
			if len(body) > 0 && body[len(body)-1] != '\n' {
				body = append(body, '\n')
			}
			body = append(body, pluginJS...)
		}
	}
	_, _ = w.Write(body)
}
