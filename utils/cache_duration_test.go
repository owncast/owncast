package utils

import "testing"

func TestCacheDurationInitSegments(t *testing.T) {
	cases := map[string]int{
		"/foo/init.mp4":            20,
		"/foo/init_0.mp4":          20,
		"/foo/init_1.mp4":          20,
		"/foo/init_10.mp4":         20,
		"/foo/offline-init.mp4":    31557600, // static embedded asset, long cache OK
		"/foo/segment.m4s":         31557600,
		"/foo/stream-foo-1.m4s":    31557600,
		"/foo/stream.m3u8":         0,
		"/foo/thumbnail.jpg":       20,
		"/foo/preview.gif":         20,
		"/foo/icon.png":            60 * 60 * 24 * 30,
		"/":                        0,
		"/foo/no-ext":              0,
		"/foo/file.html":           0,
		"/foo/font.woff2":          31557600,
		"/foo/recording-final.mp4": 31557600,
	}
	for in, want := range cases {
		if got := GetCacheDurationSecondsForPath(in); got != want {
			t.Errorf("%s: got %d, want %d", in, got, want)
		}
	}
}
