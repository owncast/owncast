package utils

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// IsHostnameInternal will attempt to determine if the hostname is internal to
// this server's network or is the loopback address.
func IsHostnameInternal(hostname string) bool {
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
