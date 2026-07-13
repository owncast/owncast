package utils

import (
	"context"
	"net"
	"net/netip"
	"os"
	"sync"
)

var (
	allowInternalFederation     bool
	allowInternalFederationOnce sync.Once
)

var nonPublicIPPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("192.88.99.0/24"),
	netip.MustParsePrefix("192.175.48.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("64:ff9b::/96"),
	netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("ff00::/8"),
}

// AllowInternalFederation returns true if the OWNCAST_ALLOW_INTERNAL_FEDERATION
// environment variable is set to "true". This is used for testing purposes only.
func AllowInternalFederation() bool {
	allowInternalFederationOnce.Do(func() {
		allowInternalFederation = os.Getenv("OWNCAST_ALLOW_INTERNAL_FEDERATION") == "true"
	})
	return allowInternalFederation
}

// IsHostnameInternal will attempt to determine if the hostname is internal to
// this server's network or is the loopback address.
// Returns false if OWNCAST_ALLOW_INTERNAL_FEDERATION is set to "true".
func IsHostnameInternal(hostname string) bool {
	return isHostnameInternal(hostname, AllowInternalFederation())
}

func isHostnameInternal(hostname string, allowInternal bool) bool {
	if allowInternal {
		return false
	}
	if ip := net.ParseIP(hostname); ip != nil {
		return isIPAddressInternal(ip)
	}

	ips, err := net.DefaultResolver.LookupIPAddr(context.Background(), hostname)
	if err != nil || len(ips) == 0 {
		return true
	}

	for _, ip := range ips {
		if isIPAddressInternal(ip.IP) {
			return true
		}
	}

	return false
}

func isIPAddressInternal(ip net.IP) bool {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return true
	}
	addr = addr.Unmap()
	if !addr.IsGlobalUnicast() {
		return true
	}
	for _, prefix := range nonPublicIPPrefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
