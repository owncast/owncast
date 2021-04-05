package transcoder

import (
	"testing"

	"github.com/owncast/owncast/models"
)

func TestFFmpegx264Command(t *testing.T) {
	latencyLevel := models.GetLatencyLevel(3)
	codec := Libx264Codec{}

	transcoder := new(Transcoder)
	transcoder.ffmpegPath = "/fake/path/ffmpeg"
	transcoder.SetInput("fakecontent.flv")
	transcoder.SetOutputPath("fakeOutput")
	transcoder.SetIdentifier("jdofFGg")
	transcoder.SetInternalHTTPPort("8123")
	transcoder.SetCodec(codec.Name())
	transcoder.currentLatencyLevel = latencyLevel

	variant := HLSVariant{}
	variant.videoBitrate = 1200
	variant.isAudioPassthrough = true
	variant.SetVideoFramerate(30)
	variant.SetCPUUsageLevel(2)
	transcoder.AddVariant(variant)

	variant2 := HLSVariant{}
	variant2.videoBitrate = 3500
	variant2.isAudioPassthrough = true
	variant2.SetVideoFramerate(24)
	variant2.SetCPUUsageLevel(4)
	transcoder.AddVariant(variant2)

	variant3 := HLSVariant{}
	variant3.isAudioPassthrough = true
	variant3.isVideoPassthrough = true
	transcoder.AddVariant(variant3)

	cmd := transcoder.getString()

	expected := `FFREPORT=file="transcoder.log":level=32 /fake/path/ffmpeg -hide_banner -loglevel warning  -i  fakecontent.flv  -map v:0 -c:v:0 libx264 -b:v:0 1200k -maxrate:v:0 1272k -g:v:0 90 -keyint_min:v:0 90 -r:v:0 30 -x264-params:v:0 "scenecut=0:open_gop=0" -bufsize:v:0 1440k -profile:v:0 high -map a:0? -c:a:0 copy -preset veryfast -map v:0 -c:v:1 libx264 -b:v:1 3500k -maxrate:v:1 3710k -g:v:1 72 -keyint_min:v:1 72 -r:v:1 24 -x264-params:v:1 "scenecut=0:open_gop=0" -bufsize:v:1 4200k -profile:v:1 high -map a:0? -c:a:1 copy -preset fast -map v:0 -c:v:2 copy -map a:0? -c:a:2 copy -preset ultrafast  -var_stream_map "v:0,a:0 v:1,a:1 v:2,a:2 " -f hls -hls_time 3 -hls_list_size 3  -tune zerolatency -pix_fmt yuv420p -sc_threshold 0 -master_pl_name stream.m3u8 -strftime 1 -hls_segment_filename http://127.0.0.1:8123/%v/stream-jdofFGg%s.ts -max_muxing_queue_size 400 -method PUT -http_persistent 0 -fflags +genpts http://127.0.0.1:8123/%v/stream.m3u8`

	if cmd != expected {
		t.Errorf("ffmpeg command does not match expected.\nGot %s\n, want: %s", cmd, expected)
	}
}
