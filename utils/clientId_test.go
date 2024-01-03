package utils

import (
	"net/http/httptest"
	"testing"
)

func TestGetIPAddressFromRequest(t *testing.T) {

	expectedIpv4Result := "203.0.113.195"
	expectedIpv6Result := "2001:9e8:5221:4201:5221:7dff:fe53"

	request := httptest.NewRequest("GET", "/", nil)

	//Test ipv4 address without X-FORWARDED-FOR header (With Port set up)
	request.RemoteAddr = "203.0.113.195:41237"
	result := GetIPAddressFromRequest(request)
	if result != expectedIpv4Result {
		t.Errorf("IPV4 Remote address only (with Port) test failed: Expected %s, got %s", expectedIpv4Result, result)
	}

	//Test ipv4 without X-FORWARDED-FOR header (Without Port set up)
	request.RemoteAddr = "203.0.113.195"
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv4Result {
		t.Errorf("IPV4 Remote address only (without Port) test failed: Expected %s, got %s", expectedIpv4Result, result)
	}

	//Test ipv4 with X-FORWARDED-FOR header (With Port SetUp)
	request.Header.Set("X-FORWARDED-FOR", "203.0.113.195:41237, 198.51.100.100:38523, 140.248.67.176:12345")
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv4Result {
		t.Errorf("IPV4 X-FORWARDED-FOR header (With Port) test failed: Expected %s, got %s", expectedIpv4Result, result)
	}

	//Test ipv4 with X-FORWARDED-FOR header (Without Port SetUp)
	request.Header.Set("X-FORWARDED-FOR", "203.0.113.195, 198.51.100.100, 140.248.67.176")
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv4Result {
		t.Errorf("IPV4 X-FORWARDED-FOR header (Without Port) test failed: Expected %s, got %s", expectedIpv4Result, result)
	}

	//Reset X-FORWARDED-FOR Header
	request.Header.Set("X-FORWARDED-FOR", "")

	//Test ipv6 address without X-FORWARDED-FOR header (With Port Set Up)
	request.RemoteAddr = "[2001:9e8:5221:4201:5221:7dff:fe53]:9080"
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv6Result {
		t.Errorf("IPV6 Remote address only (with Port) test failed: Expected %s, got %s", expectedIpv6Result, result)
	}

	//Test ipv6 address without X-FORWARDED-FOR header (Without Port Set Up)
	request.RemoteAddr = "2001:9e8:5221:4201:5221:7dff:fe53"
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv6Result {
		t.Errorf("IPV6 Remote address only (with Port) test failed: Expected %s, got %s", expectedIpv6Result, result)
	}

	//Test ipv6 with X-FORWARDED-FOR header (Without Port SetUp)
	request.Header.Set("X-FORWARDED-FOR", "2001:9e8:5221:4201:5221:7dff:fe53, 2002:9e7:5231:4211:5211:7eff:ee53, 2003:7e8:5221:4231:5221:7dff:ff53")
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv6Result {
		t.Errorf("IPV6 X-FORWARDED-FOR header (Without Port) test failed: Expected %s, got %s", expectedIpv6Result, result)
	}

	//Test ipv6 with X-FORWARDED-FOR header (With Port SetUp)
	request.Header.Set("X-FORWARDED-FOR", "[2001:9e8:5221:4201:5221:7dff:fe53]:9080, [2002:9e7:5231:4211:5211:7eff:ee53]:9080, [2003:7e8:5221:4231:5221:7dff:ff53]:9080")
	result = GetIPAddressFromRequest(request)
	if result != expectedIpv6Result {
		t.Errorf("IPV6 X-FORWARDED-FOR header (Without Port) test failed: Expected %s, got %s", expectedIpv6Result, result)
	}
}

func TestIsIpv6(t *testing.T) {
	ipv6AddressWithoutPort := "2001:9e8:5221:4201:5221:7dff:fe53"
	ipv6AddressWithPort := "[2001:9e8:5221:4201:5221:7dff:fe53]:9080"
	ipv4AddressWithoutPort := "203.0.113.195"
	ipv4AddressWithPort := "203.0.113.195:9080"

	if !isIPv6(ipv6AddressWithPort) || !isIPv6(ipv6AddressWithoutPort) {
		t.Error("IPV6 format test returned false on a valid ipv6 address")
	}

	if isIPv6(ipv4AddressWithPort) || isIPv6(ipv4AddressWithoutPort) {
		t.Error("IPV6 format test returned true on an invalid ipv6 address")
	}
}
