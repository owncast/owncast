package transcoder

import "testing"

// TestSurfaceableTranscoderMessage proves the stderr classification contract
// under ffmpeg's `-loglevel level+info`: diagnostic levels are suppressed,
// problem levels are surfaced with only the level tag (and its trailing
// space) removed so errorMap/ignoredErrors substring matching still sees the
// original message text, and rare untagged lines pass through verbatim.
func TestSurfaceableTranscoderMessage(t *testing.T) {
	cases := []struct {
		name   string
		line   string
		want   string
		wantOK bool
	}{
		{
			"plain info suppressed",
			"[info] Stream mapping:",
			"", false,
		},
		{
			"component-prefixed info suppressed",
			"[hls @ 0x55d3] [info] Opening '/x/stream0.ts' for writing",
			"", false,
		},
		{
			"warning surfaced with tag stripped",
			"[warning] VBV underflow (frame 12, -224 bits)",
			"VBV underflow (frame 12, -224 bits)", true,
		},
		{
			"component-prefixed error keeps component tag",
			"[flv @ 0x1] [error] Failed to open file 'http://127.0.0.1'",
			"[flv @ 0x1] Failed to open file 'http://127.0.0.1'", true,
		},
		{
			"untagged line surfaced verbatim",
			"Press [q] to stop, [?] for help",
			"Press [q] to stop, [?] for help", true,
		},
		{
			"fatal surfaced stripped",
			"[fatal] Unable to find a suitable output format for 'pipe:'",
			"Unable to find a suitable output format for 'pipe:'", true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := surfaceableTranscoderMessage(tc.line)
			if got != tc.want || ok != tc.wantOK {
				t.Errorf("surfaceableTranscoderMessage(%q) = (%q, %v), want (%q, %v)", tc.line, got, ok, tc.want, tc.wantOK)
			}
		})
	}
}
