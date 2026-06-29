package transcoder

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

func newTestTranscoder() *Transcoder {
	tr := new(Transcoder)
	tr.cfg = &config.Config{LogDirectory: "data/logs"}
	tr.ffmpegPath = filepath.Join("fake", "path", "ffmpeg")
	tr.SetInput("fake.flv")
	tr.SetOutputPath("fakeOutput")
	tr.SetIdentifier("test")
	tr.SetInternalHTTPPort("8123")
	codec := Libx264Codec{}
	tr.SetCodec(codec.Name())
	tr.currentLatencyLevel = models.GetLatencyLevel(2)

	v := HLSVariant{}
	v.videoBitrate = 1200
	v.isAudioPassthrough = true
	v.SetVideoFramerate(30)
	v.SetCPUUsageLevel(2)
	tr.AddVariant(v)
	return tr
}

// TestTranscoderListenerHost proves the ffmpeg HLS PUT target follows the
// configured host: the historical 127.0.0.1 default when unset, an explicit
// address otherwise, and correctly bracketed IPv6. This is the knob the
// push-to-core (remote engine) topology needs.
func TestTranscoderListenerHost(t *testing.T) {
	cases := []struct {
		name string
		host string
		want string
	}{
		{"default when unset", "", "http://127.0.0.1:8123"},
		{"custom ipv4", "10.0.0.5", "http://10.0.0.5:8123"},
		{"ipv6 bracketed", "::1", "http://[::1]:8123"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := newTestTranscoder()
			if tc.host != "" {
				tr.SetInternalHTTPHost(tc.host)
			}
			cmd := tr.GetString()
			if !strings.Contains(cmd, tc.want) {
				t.Errorf("ffmpeg command missing PUT target %q.\ncommand: %s", tc.want, cmd)
			}
		})
	}
}
