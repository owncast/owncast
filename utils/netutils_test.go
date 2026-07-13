package utils

import (
	"net"
	"testing"
)

func TestIPAddressInternal(t *testing.T) {
	tests := []struct {
		address  string
		internal bool
	}{
		{address: "0.0.0.0", internal: true},
		{address: "10.0.0.1", internal: true},
		{address: "100.64.0.1", internal: true},
		{address: "127.0.0.1", internal: true},
		{address: "169.254.169.254", internal: true},
		{address: "192.168.1.1", internal: true},
		{address: "198.18.0.1", internal: true},
		{address: "224.0.0.1", internal: true},
		{address: "::", internal: true},
		{address: "::1", internal: true},
		{address: "64:ff9b::a00:1", internal: true},
		{address: "2002:a00:1::", internal: true},
		{address: "fe80::1", internal: true},
		{address: "93.184.216.34", internal: false},
		{address: "2606:2800:220:1:248:1893:25c8:1946", internal: false},
	}

	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			got := isIPAddressInternal(net.ParseIP(test.address))
			if got != test.internal {
				t.Fatalf("isIPAddressInternal(%s) = %v, want %v", test.address, got, test.internal)
			}
		})
	}
}

func TestIsHostnameInternalFailsClosed(t *testing.T) {
	if !isHostnameInternal("does-not-exist.invalid", false) {
		t.Fatal("unresolvable hostname should be rejected")
	}
}
