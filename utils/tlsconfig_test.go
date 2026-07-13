package utils

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
)

func TestRestrictedDialerRejectsNonPublicResolvedAddress(t *testing.T) {
	dialCalled := false
	dialer := &restrictedDialer{
		lookupIP: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{
				{IP: net.ParseIP("93.184.216.34")},
				{IP: net.ParseIP("169.254.169.254")},
			}, nil
		},
		dial: func(context.Context, string, string) (net.Conn, error) {
			dialCalled = true
			return nil, errors.New("unexpected dial")
		},
	}

	_, err := dialer.DialContext(context.Background(), "tcp", "example.com:443")
	if err == nil {
		t.Fatal("expected non-public resolved address to be rejected")
	}
	if dialCalled {
		t.Fatal("dial should not run after a non-public address is resolved")
	}
}

func TestRestrictedDialerDialsValidatedIPAddress(t *testing.T) {
	var dialedAddress string
	dialer := &restrictedDialer{
		lookupIP: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}, nil
		},
		dial: func(_ context.Context, _ string, address string) (net.Conn, error) {
			dialedAddress = address
			client, server := net.Pipe()
			server.Close()
			return client, nil
		},
	}

	conn, err := dialer.DialContext(context.Background(), "tcp", "example.com:443")
	if err != nil {
		t.Fatalf("expected public address to be dialed: %v", err)
	}
	conn.Close()
	if dialedAddress != "93.184.216.34:443" {
		t.Fatalf("dialed %q instead of the validated IP address", dialedAddress)
	}
}

func TestOutboundRedirectPolicyRejectsNonPublicHost(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://169.254.169.254/latest", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := outboundRedirectPolicy(false)(req, nil); err == nil {
		t.Fatal("expected redirect to link-local address to be rejected")
	}
}

func TestValidatePublicHTTPSURL(t *testing.T) {
	if err := ValidatePublicHTTPSURL("https://93.184.216.34/push"); err != nil {
		t.Fatalf("expected public HTTPS URL to be accepted: %v", err)
	}
	for _, rawURL := range []string{
		"http://93.184.216.34/push",
		"https://127.0.0.1/push",
		"https://100.64.0.1/push",
		"https://[64:ff9b::a00:1]/push",
	} {
		t.Run(rawURL, func(t *testing.T) {
			if err := ValidatePublicHTTPSURL(rawURL); err == nil {
				t.Fatalf("expected %q to be rejected", rawURL)
			}
		})
	}
}
