package authrepository

import (
	"testing"

	"github.com/owncast/owncast/services/datastore"
)

func TestBanIPAddressValidatesAddress(t *testing.T) {
	ds, err := datastore.SetupPersistence(":memory:", t.TempDir())
	if err != nil {
		t.Fatalf("SetupPersistence: %v", err)
	}
	defer ds.DB.Close() //nolint

	repository := New(ds)

	for _, address := range []string{"not-an-ip", "", "127.0.0.1:8080"} {
		t.Run("rejects "+address, func(t *testing.T) {
			if err := repository.BanIPAddress(address, "test"); err == nil {
				t.Errorf("BanIPAddress(%q) succeeded, want error", address)
			}
		})
	}

	for _, address := range []string{"127.0.0.1", "2001:db8::1"} {
		t.Run("accepts "+address, func(t *testing.T) {
			if err := repository.BanIPAddress(address, "test"); err != nil {
				t.Errorf("BanIPAddress(%q): %v", address, err)
			}
		})
	}
}
