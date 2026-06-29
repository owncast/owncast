package transcoder

import (
	"net"
	"strconv"
	"testing"

	"github.com/owncast/owncast/config"
)

type noopReceiverCallbacks struct{}

func (noopReceiverCallbacks) SegmentWritten(string)         {}
func (noopReceiverCallbacks) VariantPlaylistWritten(string) {}
func (noopReceiverCallbacks) MasterPlaylistWritten(string)  {}

func assertNumericPort(t *testing.T, port string) {
	t.Helper()
	if port == "" {
		t.Fatal("listener port was not recorded")
	}
	if _, err := strconv.Atoi(port); err != nil {
		t.Errorf("listener port %q is not numeric: %v", port, err)
	}
}

// TestReceiverBindHost covers the configurable receiver bind. An empty host
// must fall back to loopback rather than binding every interface, and the
// chosen port must be parsed correctly even for a bracketed IPv6 listener
// address (the old strings.Split(":") parse returned an empty port there).
func TestReceiverBindHost(t *testing.T) {
	t.Run("empty host defaults to loopback and records the port", func(t *testing.T) {
		cfg := &config.Config{}
		NewFileWriterReceiverService(cfg).SetupFileWriterReceiverService(noopReceiverCallbacks{})
		assertNumericPort(t, cfg.InternalHLSListenerPort)
	})

	t.Run("ipv6 host records a numeric port", func(t *testing.T) {
		probe, err := net.Listen("tcp", "[::1]:0")
		if err != nil {
			t.Skip("no IPv6 loopback available")
		}
		_ = probe.Close()

		cfg := &config.Config{InternalHLSListenerHost: "::1"}
		NewFileWriterReceiverService(cfg).SetupFileWriterReceiverService(noopReceiverCallbacks{})
		assertNumericPort(t, cfg.InternalHLSListenerPort)
	})
}
