package transcoder

import "testing"

func TestIsCodecSupported(t *testing.T) {
	for _, codec := range []string{
		"libx264",
		"h264_omx",
		"h264_vaapi",
		"h264_qsv",
		"h264_nvenc",
		"h264_v4l2m2m",
		"h264_videotoolbox",
	} {
		if !IsCodecSupported(codec) {
			t.Errorf("expected %q to be supported", codec)
		}
	}
	if IsCodecSupported("h264") {
		t.Fatal("expected unknown codec to be rejected")
	}
}
