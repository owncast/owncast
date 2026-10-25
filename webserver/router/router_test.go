package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type deadlineResponseRecorder struct {
	*httptest.ResponseRecorder
	readDeadline time.Time
}

func (r *deadlineResponseRecorder) SetReadDeadline(deadline time.Time) error {
	r.readDeadline = deadline
	return nil
}

func TestPluginUploadReadTimeout(t *testing.T) {
	handler := withPluginUploadReadTimeout(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	tests := []struct {
		name             string
		method           string
		path             string
		wantReadDeadline bool
	}{
		{name: "plugin upload", method: http.MethodPost, path: "/api/admin/plugins", wantReadDeadline: true},
		{name: "plugin upload subpath", method: http.MethodPost, path: "/api/admin/plugins/example", wantReadDeadline: false},
		{name: "plugin list", method: http.MethodGet, path: "/api/admin/plugins", wantReadDeadline: false},
		{name: "other API", method: http.MethodPost, path: "/api/admin/config", wantReadDeadline: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := &deadlineResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
			before := time.Now()
			handler.ServeHTTP(writer, httptest.NewRequest(test.method, test.path, nil))

			if got := !writer.readDeadline.IsZero(); got != test.wantReadDeadline {
				t.Fatalf("read deadline set = %t, want %t", got, test.wantReadDeadline)
			}
			if test.wantReadDeadline && writer.readDeadline.Before(before.Add(pluginUploadReadTimeout-time.Second)) {
				t.Fatalf("read deadline = %v, want approximately %v from now", writer.readDeadline, pluginUploadReadTimeout)
			}
		})
	}
}
