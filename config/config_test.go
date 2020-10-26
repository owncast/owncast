package config

import "testing"

func TestDefaults(t *testing.T) {
	_default = getDefaults()

	encoderPreset := "veryfast"
	framerate := 24

	quality := StreamQuality{}
	if quality.GetEncoderPreset() != encoderPreset {
		t.Errorf("default encoder preset does not match expected.  Got %s, want: %s", quality.GetEncoderPreset(), encoderPreset)
	}

	if quality.GetFramerate() != framerate {
		t.Errorf("default framerate does not match expected.  Got %d, want: %d", quality.GetFramerate(), framerate)
	}
}
