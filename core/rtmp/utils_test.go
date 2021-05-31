package rtmp

import "testing"

func Test_secretMatch(t *testing.T) {
	tests := []struct {
		name      string
		streamKey string
		url       string
		want      bool
	}{
		{"simple positive", "abc", "rtmp://example.com/live/abc", true},
		{"simple negative", "abc", "rtmp://example.com/live/def", false},
		{"simple positive with numbers", "abc123", "rtmp://example.com/live/abc123", true},
		{"simple negative with numbers", "abc123", "rtmp://example.com/live/def456", false},
		{"positive with url chars", "one/two/three", "rtmp://example.com/live/one/two/three", true},
		{"negative with url chars", "one/two/three", "rtmp://example.com/live/four/five/six", false},
		{"check the entire secret", "three", "rtmp://example.com/live/one/two/three", false},
		{"with /live/ in secret", "one/live/three", "rtmp://example.com/live/one/live/three", true},
		{"bad urls", "anything", "nonsense", false},
		{"missing secret", "abc", "rtmp://example.com/live/", false},
		{"missing secret and missing last slash", "abc", "rtmp://example.com/live", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := secretMatch(tt.streamKey, tt.url); got != tt.want {
				t.Errorf("secretMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
