package utils

import (
	"net"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	allowInternalFederation     bool
	allowInternalFederationOnce sync.Once
)

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
	// Allow internal federation for testing purposes.
	if AllowInternalFederation() {
		return false
	}
	// If this is already an IP address don't try to resolve it
	if ip := net.ParseIP(hostname); ip != nil {
		return isIPAddressInternal(ip)
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		// Default to false if we can't resolve the hostname.
		log.Debugln("Unable to resolve hostname:", hostname)
		return false
	}

	for _, ip := range ips {
		if isIPAddressInternal(ip) {
			return true
		}
	}

	return false
}

func isIPAddressInternal(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate()
}
