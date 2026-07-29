package router

import (
	"crypto/tls"
	"errors"
	"net"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/acme"
)

func TestBindAutoTLSListenersAllOrNothing(t *testing.T) {
	// Occupy one port so the pair bind fails partway through.
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	occupiedPort := occupied.Addr().(*net.TCPAddr).Port

	// Find a free port for the http side.
	free, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	freePort := free.Addr().(*net.TCPAddr).Port
	free.Close()

	httpLn, httpsLn, err := bindAutoTLSListeners("127.0.0.1", freePort, occupiedPort)
	if err == nil {
		httpLn.Close()
		httpsLn.Close()
		t.Fatal("expected an error when the https port is occupied")
	}

	// The successfully-bound http listener must have been released.
	relisten, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(freePort)))
	if err != nil {
		t.Fatalf("http listener was not released after partial bind failure: %v", err)
	}
	relisten.Close()
}

func TestBindAutoTLSListenersSuccess(t *testing.T) {
	a, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	portA := a.Addr().(*net.TCPAddr).Port
	a.Close()
	b, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	portB := b.Addr().(*net.TCPAddr).Port
	b.Close()

	httpLn, httpsLn, err := bindAutoTLSListeners("127.0.0.1", portA, portB)
	if err != nil {
		t.Fatalf("expected both listeners to bind: %v", err)
	}
	httpLn.Close()
	httpsLn.Close()
}

func TestGetCertificateBackoff(t *testing.T) {
	calls := 0
	at := &autoTLS{host: "live.example.com"}
	at.lookup = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		calls++
		return nil, errors.New("issuance failed")
	}

	hello := &tls.ClientHelloInfo{ServerName: "live.example.com"}

	// First failure reaches the lookup and arms the backoff.
	if _, err := at.getCertificate(hello); err == nil {
		t.Fatal("expected a failure")
	}
	if calls != 1 {
		t.Fatalf("expected 1 lookup call, got %d", calls)
	}

	// While backing off, handshakes fail without another lookup.
	if _, err := at.getCertificate(hello); err == nil {
		t.Fatal("expected a backoff error")
	}
	if calls != 1 {
		t.Fatalf("backoff must prevent further lookups, got %d calls", calls)
	}

	// After the window elapses the lookup is retried.
	at.mu.Lock()
	at.lastFailure = time.Now().Add(-acmeFailureBackoff - time.Second)
	at.mu.Unlock()
	if _, err := at.getCertificate(hello); err == nil {
		t.Fatal("expected a failure")
	}
	if calls != 2 {
		t.Fatalf("expected a retry after the backoff window, got %d calls", calls)
	}
}

func TestGetCertificateSuccessResetsBackoff(t *testing.T) {
	cert := &tls.Certificate{}
	fail := true
	at := &autoTLS{host: "live.example.com"}
	at.lookup = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		if fail {
			return nil, errors.New("issuance failed")
		}
		return cert, nil
	}
	hello := &tls.ClientHelloInfo{ServerName: "live.example.com"}

	if _, err := at.getCertificate(hello); err == nil {
		t.Fatal("expected a failure")
	}

	// Simulate the window elapsing, then succeed.
	at.mu.Lock()
	at.lastFailure = time.Now().Add(-acmeFailureBackoff - time.Second)
	at.mu.Unlock()
	fail = false
	got, err := at.getCertificate(hello)
	if err != nil || got != cert {
		t.Fatalf("expected the certificate, got %v / %v", got, err)
	}

	// Backoff must be cleared after success.
	at.mu.Lock()
	cleared := at.lastFailure.IsZero()
	at.mu.Unlock()
	if !cleared {
		t.Fatal("a successful lookup must clear the failure backoff")
	}
}

func TestGetCertificateOtherHostsDoNotArmBackoff(t *testing.T) {
	at := &autoTLS{host: "live.example.com"}
	at.lookup = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return nil, errors.New("host not configured")
	}

	// A scanner handshake with a foreign SNI fails but must not arm the
	// backoff for the real hostname.
	if _, err := at.getCertificate(&tls.ClientHelloInfo{ServerName: "scanner.example.net"}); err == nil {
		t.Fatal("expected a failure for a foreign host")
	}
	at.mu.Lock()
	armed := !at.lastFailure.IsZero()
	at.mu.Unlock()
	if armed {
		t.Fatal("a foreign-host failure must not arm the backoff")
	}
}

func TestGetCertificateACMEChallengeBypassesBackoff(t *testing.T) {
	challengeCert := &tls.Certificate{}
	calls := 0
	at := &autoTLS{
		host:        "live.example.com",
		lastFailure: time.Now(),
		lookup: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			calls++
			return challengeCert, nil
		},
	}

	hello := &tls.ClientHelloInfo{
		ServerName:      "live.example.com",
		SupportedProtos: []string{acme.ALPNProto},
	}
	got, err := at.getCertificate(hello)
	if err != nil || got != challengeCert {
		t.Fatalf("ACME challenge lookup must bypass issuance backoff, got %v / %v", got, err)
	}
	if calls != 1 {
		t.Fatalf("expected the ACME challenge lookup to run once, got %d calls", calls)
	}
}

func TestGetCertificateACMEChallengeFailureDoesNotArmBackoff(t *testing.T) {
	at := &autoTLS{
		host: "live.example.com",
		lookup: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			return nil, errors.New("challenge certificate unavailable")
		},
	}

	hello := &tls.ClientHelloInfo{
		ServerName:      "live.example.com",
		SupportedProtos: []string{acme.ALPNProto},
	}
	if _, err := at.getCertificate(hello); err == nil {
		t.Fatal("expected the challenge lookup to fail")
	}
	at.mu.Lock()
	armed := !at.lastFailure.IsZero()
	at.mu.Unlock()
	if armed {
		t.Fatal("an ACME challenge lookup failure must not arm issuance backoff")
	}
}

func TestGetCertificateMatchesTrailingDotSNI(t *testing.T) {
	calls := 0
	at := &autoTLS{host: "live.example.com"}
	at.lookup = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		calls++
		return nil, errors.New("issuance failed")
	}

	// A fully-qualified SNI with a trailing dot is the configured host.
	if _, err := at.getCertificate(&tls.ClientHelloInfo{ServerName: "live.example.com."}); err == nil {
		t.Fatal("expected a failure")
	}
	at.mu.Lock()
	armed := !at.lastFailure.IsZero()
	at.mu.Unlock()
	if !armed {
		t.Fatal("a trailing-dot SNI for the configured host must arm the backoff")
	}
}

func TestEnvPortOr(t *testing.T) {
	t.Setenv("OWNCAST_TLS_HTTP_PORT", "")
	if got := envPortOr("OWNCAST_TLS_HTTP_PORT", 80); got != 80 {
		t.Fatalf("empty value must return the default, got %d", got)
	}
	t.Setenv("OWNCAST_TLS_HTTP_PORT", "5080")
	if got := envPortOr("OWNCAST_TLS_HTTP_PORT", 80); got != 5080 {
		t.Fatalf("expected the override, got %d", got)
	}
	t.Setenv("OWNCAST_TLS_HTTP_PORT", "not-a-port")
	if got := envPortOr("OWNCAST_TLS_HTTP_PORT", 80); got != 80 {
		t.Fatalf("invalid value must return the default, got %d", got)
	}
	t.Setenv("OWNCAST_TLS_HTTP_PORT", "70000")
	if got := envPortOr("OWNCAST_TLS_HTTP_PORT", 80); got != 80 {
		t.Fatalf("out-of-range value must return the default, got %d", got)
	}
}
