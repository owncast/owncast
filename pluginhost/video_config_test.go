package pluginhost

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/plugins"
)

func newVideoConfigEnv(t *testing.T) (*plugins.HostEnv, configrepository.ConfigRepository) {
	t.Helper()
	ds, err := datastore.SetupPersistence(":memory:", t.TempDir())
	if err != nil {
		t.Fatalf("setup persistence: %v", err)
	}
	repo := configrepository.New(ds)
	env := &plugins.HostEnv{}
	wireVideoConfigHostFns(env, Deps{ConfigRepository: repo})
	return env, repo
}

func TestVideoConfigReadMatchesAdminSettings(t *testing.T) {
	env, repo := newVideoConfigEnv(t)
	if err := repo.SetAutoplay(models.AutoplaySoundOnly); err != nil {
		t.Fatal(err)
	}
	if err := repo.SetVideoCodec("h264_vaapi"); err != nil {
		t.Fatal(err)
	}
	variants := []models.StreamOutputVariant{{
		ScaledWidth:   1280,
		ScaledHeight:  720,
		Framerate:     30,
		VideoBitrate:  3000,
		AudioBitrate:  128,
		CPUUsageLevel: 3,
	}}
	if err := repo.SetStreamOutputVariants(variants); err != nil {
		t.Fatal(err)
	}

	got := env.VideoConfig()
	want := plugins.VideoConfig{
		LatencyLevel: repo.GetStreamLatencyLevel().Level,
		Codec:        "h264_vaapi",
		Autoplay:     models.AutoplaySoundOnly,
		Variants: []plugins.StreamVariant{{
			Width:         1280,
			Height:        720,
			Framerate:     30,
			VideoBitrate:  3000,
			CPUUsageLevel: 3,
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("video config = %#v, want %#v", got, want)
	}

	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "audioBitrate") {
		t.Fatalf("plugin video config unexpectedly exposes audio bitrate: %s", data)
	}
}

func TestVideoConfigWriteMatchesAdminSettings(t *testing.T) {
	env, repo := newVideoConfigEnv(t)
	if err := repo.SetStreamOutputVariants([]models.StreamOutputVariant{{
		Name:               "Existing output",
		AudioBitrate:       160,
		IsAudioPassthrough: false,
	}}); err != nil {
		t.Fatal(err)
	}
	autoplay := models.AutoplayAlways
	codec := "h264_nvenc"
	update := plugins.VideoConfigUpdate{
		Autoplay: &autoplay,
		Codec:    &codec,
		Variants: []plugins.StreamVariant{{
			Width:         1920,
			Height:        1080,
			Framerate:     30,
			VideoBitrate:  6000,
			CPUUsageLevel: 1,
		}},
	}
	if err := env.WriteVideoConfig("video-settings", update); err != nil {
		t.Fatal(err)
	}
	if got := repo.GetAutoplay(); got != autoplay {
		t.Fatalf("autoplay = %q, want %q", got, autoplay)
	}
	if got := repo.GetVideoCodec(); got != codec {
		t.Fatalf("codec = %q, want %q", got, codec)
	}
	variants := repo.GetStreamOutputVariants()
	if len(variants) != 1 {
		t.Fatalf("variant count = %d, want 1", len(variants))
	}
	if got := variants[0].CPUUsageLevel; got != 1 {
		t.Fatalf("CPU usage level = %d, want 1", got)
	}
	if got := variants[0].AudioBitrate; got != 160 {
		t.Fatalf("audio bitrate = %d, want preserved value 160", got)
	}
	if got := variants[0].Name; got != "Existing output" {
		t.Fatalf("name = %q, want preserved internal name", got)
	}
}

func TestVideoConfigWriteRejectsInvalidEnums(t *testing.T) {
	env, _ := newVideoConfigEnv(t)

	invalidAutoplay := models.AutoplayMode("sometimes")
	if err := env.WriteVideoConfig("video-settings", plugins.VideoConfigUpdate{Autoplay: &invalidAutoplay}); err == nil {
		t.Fatal("expected invalid autoplay to be rejected")
	}

	invalidCodec := "h264"
	if err := env.WriteVideoConfig("video-settings", plugins.VideoConfigUpdate{Codec: &invalidCodec}); err == nil {
		t.Fatal("expected invalid codec to be rejected")
	}
}
