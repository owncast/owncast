package ffmpeg

import (
	"testing"
)

func TestFFmpegCommand(t *testing.T) {
	transcoder := new(Transcoder)
	transcoder.ffmpegPath = "/fake/path/ffmpeg"
	transcoder.SetSegmentLength(4)
	transcoder.SetInput("fakecontent.flv")
	transcoder.SetOutputPath("fakeOutput")
	transcoder.SetHLSPlaylistLength(10)
	transcoder.SetIdentifier("jdofFGg")
	transcoder.SetInternalHTTPPort(8123)

	variant := HLSVariant{}
	variant.videoBitrate = 1200
	variant.isAudioPassthrough = true
	variant.encoderPreset = "veryfast"
	variant.SetVideoFramerate(30)

	transcoder.AddVariant(variant)

	cmd := transcoder.getString()

	expected := `/fake/path/ffmpeg -hide_banner -loglevel warning -i  fakecontent.flv  -map v:0 -c:v:0 libx264 -b:v:0 1200k -maxrate:v:0 1272k -bufsize:v:0 1440k -g:v:0 119 -x264-params:v:0 "scenecut=0:open_gop=0:min-keyint=119:keyint=119" -map a:0 -c:a:0 copy -r 30 -preset veryfast  -var_stream_map "v:0,a:0 " -f hls -hls_time 4 -hls_list_size 10 -hls_delete_threshold 10  -tune zerolatency -pix_fmt yuv420p -profile:v high -sc_threshold 0 -master_pl_name stream.m3u8 -strftime 1 -hls_segment_filename http://127.0.0.1:8123/%v/stream-jdofFGg%s.ts -max_muxing_queue_size 400 -method PUT -http_persistent 1 -fflags +genpts http://127.0.0.1:8123/%v/stream.m3u8 2> transcoder.log`

	if cmd != expected {
		t.Errorf("ffmpeg command does not match expected.  Got %s, want: %s", cmd, expected)
	}
}
