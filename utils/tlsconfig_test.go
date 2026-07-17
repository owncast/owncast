package utils

import (
	"context"
	"errors"
	"net"
	"net/http"
	"slices"
	"testing"
	"time"
)

func TestRestrictedDialerSkipsNonPublicResolvedAddresses(t *testing.T) {
	var dialed []string
	dialer := &restrictedDialer{
		lookupIP: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{
				{IP: net.ParseIP("169.254.169.254")},
				{IP: net.ParseIP("93.184.216.34")},
			}, nil
		},
		dial: func(_ context.Context, _ string, address string) (net.Conn, error) {
			dialed = append(dialed, address)
			client, server := net.Pipe()
			server.Close()
			return client, nil
		},
	}

	conn, err := dialer.DialContext(context.Background(), "tcp", "example.com:443")
	if err != nil {
		t.Fatalf("expected the public address to be dialed: %v", err)
	}
	conn.Close()
	if !slices.Equal(dialed, []string{"93.184.216.34:443"}) {
		t.Fatalf("dialed %v, want only the public address", dialed)
	}
}

func TestRestrictedDialerRejectsHostsWithOnlyNonPublicAddresses(t *testing.T) {
	dialCalled := false
	dialer := &restrictedDialer{
		lookupIP: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{
				{IP: net.ParseIP("169.254.169.254")},
				{IP: net.ParseIP("10.0.0.1")},
			}, nil
		},
		dial: func(context.Context, string, string) (net.Conn, error) {
			dialCalled = true
			return nil, errors.New("unexpected dial")
		},
	}

	_, err := dialer.DialContext(context.Background(), "tcp", "example.com:443")
	if err == nil {
		t.Fatal("expected host with no public addresses to be rejected")
	}
	if dialCalled {
		t.Fatal("dial should not run when no public address is resolved")
	}
}

func TestRestrictedDialerInterleavesAddressFamilies(t *testing.T) {
	var dialed []string
	dialer := &restrictedDialer{
		lookupIP: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{
				{IP: net.ParseIP("2606:2800:220:1::1")},
				{IP: net.ParseIP("2606:2800:220:1::2")},
				{IP: net.ParseIP("93.184.216.34")},
			}, nil
		},
		dial: func(_ context.Context, _ string, address string) (net.Conn, error) {
			dialed = append(dialed, address)
			return nil, errors.New("dial failed")
		},
	}

	if _, err := dialer.DialContext(context.Background(), "tcp", "example.com:443"); err == nil {
		t.Fatal("expected dial failure")
	}
	want := []string{"[2606:2800:220:1::1]:443", "93.184.216.34:443", "[2606:2800:220:1::2]:443"}
	if !slices.Equal(dialed, want) {
		t.Fatalf("dial order %v, want %v", dialed, want)
	}
}

func TestRestrictedDialerSplitsDeadlineAcrossAttempts(t *testing.T) {
	var deadlines []time.Time
	dialer := &restrictedDialer{
		lookupIP: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{
				{IP: net.ParseIP("93.184.216.34")},
				{IP: net.ParseIP("93.184.216.35")},
			}, nil
		},
		dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("expected each dial attempt to carry a deadline")
			}
			deadlines = append(deadlines, deadline)
			return nil, errors.New("dial failed")
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if _, err := dialer.DialContext(ctx, "tcp", "example.com:443"); err == nil {
		t.Fatal("expected dial failure")
	}
	parentDeadline, _ := ctx.Deadline()
	if len(deadlines) != 2 {
		t.Fatalf("expected 2 dial attempts, got %d", len(deadlines))
	}
	// The first attempt gets roughly half the 8s budget, not all of it.
	if !deadlines[0].Before(parentDeadline.Add(-2 * time.Second)) {
		t.Fatalf("first attempt deadline %v should be well before the parent deadline %v", deadlines[0], parentDeadline)
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

func TestOutboundRedirectPolicyRejectsHTTPSDowngrade(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://93.184.216.34/push", nil)
	if err != nil {
		t.Fatal(err)
	}
	previous, err := http.NewRequest(http.MethodGet, "https://93.184.216.34/push", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := outboundRedirectPolicy(false)(req, []*http.Request{previous}); err == nil {
		t.Fatal("expected HTTPS to HTTP redirect to be rejected")
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
