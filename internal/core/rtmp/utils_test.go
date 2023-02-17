package rtmp

import "testing"

func Test_secretMatch(t *testing.T) {
	tests := []struct {
		name      string
		streamKey string
		path      string
		want      bool
	}{
		{"positive", "abc", "/live/abc", true},
		{"negative", "abc", "/live/def", false},
		{"positive with numbers", "abc123", "/live/abc123", true},
		{"negative with numbers", "abc123", "/live/def456", false},
		{"positive with url chars", "one/two/three", "/live/one/two/three", true},
		{"negative with url chars", "one/two/three", "/live/four/five/six", false},
		{"check the entire secret", "three", "/live/one/two/three", false},
		{"with /live/ in secret", "one/live/three", "/live/one/live/three", true},
		{"bad path", "anything", "nonsense", false},
		{"missing secret", "abc", "/live/", false},
		{"missing secret and missing last slash", "abc", "/live", false},
		{"streamkey before /live/", "streamkey", "/streamkey/live", false},
		{"missing /live/", "anything", "/something/else", false},
		{"stuff before and after /live/", "after", "/before/live/after", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := secretMatch(tt.streamKey, tt.path); got != tt.want {
				t.Errorf("secretMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
