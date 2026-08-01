package pluginhost

import (
	"testing"

	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"

	"github.com/owncast/owncast/services/plugins"
)

func TestPluginLoggingPreservesLevelAndIdentity(t *testing.T) {
	originalLevel := log.GetLevel()
	log.SetLevel(log.InfoLevel)
	defer log.SetLevel(originalLevel)

	hook := logtest.NewGlobal()
	defer hook.Reset()

	env := &plugins.HostEnv{}
	wirePluginLoggingHostFn(env)

	cases := []struct {
		level plugins.PluginLogLevel
		want  log.Level
	}{
		{plugins.PluginLogInfo, log.InfoLevel},
		{plugins.PluginLogWarning, log.WarnLevel},
		{plugins.PluginLogError, log.ErrorLevel},
	}
	for _, tc := range cases {
		env.Log("weather-alerts", tc.level, "forecast changed")
	}

	entries := hook.AllEntries()
	if len(entries) != len(cases) {
		t.Fatalf("log entries = %d, want %d", len(entries), len(cases))
	}
	for i, tc := range cases {
		if entries[i].Level != tc.want {
			t.Errorf("entry %d level = %s, want %s", i, entries[i].Level, tc.want)
		}
		if entries[i].Message != "plugin weather-alerts: forecast changed" {
			t.Errorf("entry %d message = %q", i, entries[i].Message)
		}
	}
}
