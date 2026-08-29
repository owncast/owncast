package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestAcceptsVideo(t *testing.T) {
	tests := []struct {
		name   string
		accept []string
		want   bool
	}{
		{
			name:   "video media type",
			accept: []string{"video/mp4"},
			want:   true,
		},
		{
			name:   "video wildcard with parameters",
			accept: []string{"text/html, video/*; q=0.5"},
			want:   true,
		},
		{
			name:   "HLS media type",
			accept: []string{"application/vnd.apple.mpegurl"},
			want:   true,
		},
		{
			name:   "legacy HLS media type",
			accept: []string{"application/x-mpegURL"},
			want:   true,
		},
		{
			name:   "separate video accept header",
			accept: []string{"text/html", "video/webm"},
			want:   true,
		},
		{
			name:   "unacceptable video media type",
			accept: []string{"video/mp4; q=0"},
			want:   false,
		},
		{
			name:   "invalid quality",
			accept: []string{"video/mp4; q=2"},
			want:   false,
		},
		{
			name:   "web media types",
			accept: []string{"text/html, application/xhtml+xml, image/avif"},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			for _, accept := range test.accept {
				req.Header.Add("Accept", accept)
			}

			if got := requestAcceptsVideo(req); got != test.want {
				t.Errorf("requestAcceptsVideo() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIndexHandlerRedirectsVideoRequestToHLS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "video/mp4")
	response := httptest.NewRecorder()

	(&Handlers{}).IndexHandler(response, req)

	if response.Code != http.StatusTemporaryRedirect {
		t.Errorf("status = %d, want %d", response.Code, http.StatusTemporaryRedirect)
	}
	if location := response.Header().Get("Location"); location != "/hls/stream.m3u8" {
		t.Errorf("Location = %q, want %q", location, "/hls/stream.m3u8")
	}
}
