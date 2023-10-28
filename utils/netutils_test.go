package utils

import (
	"net"
	"testing"
)

func TestIPAddressInternal(t *testing.T) {
	internalLoopbackHost := "localhost"
	internalLoopbackHostTest := IsHostnameInternal(internalLoopbackHost)
	if !internalLoopbackHostTest {
		t.Errorf("IsHostnameInternal(%s) = %v; want true", internalLoopbackHost, internalLoopbackHostTest)
	}

	internalLoopbackIP := net.ParseIP("127.0.0.1")
	internalLoopbackIPTest := isIPAddressInternal(internalLoopbackIP)
	if !internalLoopbackIPTest {
		t.Errorf("isIPAddressInternal(%s) = %v; want true", internalLoopbackIP, internalLoopbackIPTest)
	}

	externalHost := "example.com"
	externalHostTest := IsHostnameInternal(externalHost)
	if externalHostTest {
		t.Errorf("IsHostnameInternal(%s) = %v; want false", externalHost, externalHostTest)
	}

	externalIP := net.ParseIP("93.184.216.34")
	externalIPTest := isIPAddressInternal(externalIP)
	if externalIPTest {
		t.Errorf("isIPAddressInternal(%s) = %v; want false", externalIP, externalIPTest)
	}
}
