package transcoder

import (
	"path/filepath"
	"testing"

	"github.com/owncast/owncast/models"
)

func TestFFmpegVaapiCommand(t *testing.T) {
	latencyLevel := models.GetLatencyLevel(2)
	codec := VaapiCodec{}

	transcoder := new(Transcoder)
	transcoder.ffmpegPath = filepath.Join("fake", "path", "ffmpeg")
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

	expectedLogPath := filepath.Join("data", "logs", "transcoder.log")
	expected := `FFREPORT=file="` + expectedLogPath + `":level=32 ` + transcoder.ffmpegPath + ` -hide_banner -loglevel warning -vaapi_device /dev/dri/renderD128 -fflags +genpts -i  fakecontent.flv  -map v:0 -c:v:0 h264_vaapi -b:v:0 1008k -maxrate:v:0 1088k -g:v:0 90 -keyint_min:v:0 90 -r:v:0 30  -map a:0? -c:a:0 copy -filter:v:0 "format=nv12,hwupload" -preset veryfast -map v:0 -c:v:1 h264_vaapi -b:v:1 3308k -maxrate:v:1 3572k -g:v:1 72 -keyint_min:v:1 72 -r:v:1 24  -map a:0? -c:a:1 copy -filter:v:1 "format=nv12,hwupload" -preset fast -map v:0 -c:v:2 copy -map a:0? -c:a:2 copy -preset ultrafast  -var_stream_map "v:0,a:0 v:1,a:1 v:2,a:2 " -f hls -hls_time 3 -hls_list_size 10 -hls_flags program_date_time+independent_segments+omit_endlist  -segment_format_options mpegts_flags=mpegts_copyts=1  -pix_fmt vaapi_vld -sc_threshold 0 -master_pl_name stream.m3u8 -strftime 1 -hls_segment_filename http://127.0.0.1:8123/%v/stream-jdofFGg%s.ts -max_muxing_queue_size 400 -method PUT -http_persistent 0 http://127.0.0.1:8123/%v/stream.m3u8`

	if cmd != expected {
		t.Errorf("ffmpeg command does not match expected.\nGot %s\n, want: %s", cmd, expected)
	}
}
