package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

const (
	httpScheme  = "http"
	httpsScheme = "https"
)

var (
	insecureSkipVerify     bool
	insecureSkipVerifyOnce sync.Once
)

// IsInsecureSkipVerifyEnabled returns true if the OWNCAST_INSECURE_SKIP_VERIFY
// environment variable is set to "true". This is intended for testing only.
func IsInsecureSkipVerifyEnabled() bool {
	insecureSkipVerifyOnce.Do(func() {
		insecureSkipVerify = os.Getenv("OWNCAST_INSECURE_SKIP_VERIFY") == "true"
		if insecureSkipVerify {
			log.Warnln("OWNCAST_INSECURE_SKIP_VERIFY is enabled - TLS certificate verification disabled (testing only)")
		}
	})
	return insecureSkipVerify
}

// GetTLSConfig returns a TLS config that optionally skips certificate verification
// based on the OWNCAST_INSECURE_SKIP_VERIFY environment variable.
func GetTLSConfig() *tls.Config {
	if IsInsecureSkipVerifyEnabled() {
		return &tls.Config{
			InsecureSkipVerify: true, // #nosec G402 - intentional for testing
		}
	}
	return nil
}

// GetHTTPTransportWithTLS returns an http.Transport configured with TLS settings.
// If OWNCAST_INSECURE_SKIP_VERIFY is set, certificate verification is skipped.
func GetHTTPTransportWithTLS(baseTransport *http.Transport) *http.Transport {
	if baseTransport == nil {
		baseTransport = &http.Transport{}
	}
	baseTransport.TLSClientConfig = GetTLSConfig()
	return baseTransport
}

type restrictedDialer struct {
	allowInternal bool
	lookupIP      func(context.Context, string) ([]net.IPAddr, error)
	dial          func(context.Context, string, string) (net.Conn, error)
}

func newRestrictedDialer(allowInternal bool) *restrictedDialer {
	dialer := &net.Dialer{
		Timeout:   8 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &restrictedDialer{
		allowInternal: allowInternal,
		lookupIP:      net.DefaultResolver.LookupIPAddr,
		dial:          dialer.DialContext,
	}
}

func (d *restrictedDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid outbound address %q: %w", address, err)
	}

	var addresses []net.IPAddr
	if ip := net.ParseIP(host); ip != nil {
		addresses = []net.IPAddr{{IP: ip}}
	} else {
		addresses, err = d.lookupIP(ctx, host)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve outbound host %q: %w", host, err)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("outbound host %q resolved to no addresses", host)
	}

	for _, address := range addresses {
		if !d.allowInternal && isIPAddressInternal(address.IP) {
			return nil, fmt.Errorf("refusing to connect to non-public address %s for host %q", address.IP, host)
		}
	}

	var dialErrors []error
	for _, address := range addresses {
		dialAddress := net.JoinHostPort(address.IP.String(), port)
		conn, err := d.dial(ctx, network, dialAddress)
		if err == nil {
			return conn, nil
		}
		dialErrors = append(dialErrors, err)
	}

	return nil, fmt.Errorf("unable to connect to outbound host %q: %w", host, errors.Join(dialErrors...))
}

func newOutboundHTTPClient(allowInternal bool) *http.Client {
	transport := GetHTTPTransportWithTLS(&http.Transport{
		DialContext:         newRestrictedDialer(allowInternal).DialContext,
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     10 * time.Second,
	})

	return &http.Client{
		Transport:     transport,
		Timeout:       8 * time.Second,
		CheckRedirect: outboundRedirectPolicy(allowInternal),
	}
}

func outboundRedirectPolicy(allowInternal bool) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if req.URL.Scheme != httpScheme && req.URL.Scheme != httpsScheme {
			return fmt.Errorf("refusing to follow redirect using %q scheme", req.URL.Scheme)
		}
		if len(via) > 0 && via[len(via)-1].URL.Scheme == httpsScheme && req.URL.Scheme != httpsScheme {
			return errors.New("refusing to follow redirect from HTTPS to HTTP")
		}
		if isHostnameInternal(req.URL.Hostname(), allowInternal) {
			return fmt.Errorf("refusing to follow redirect to non-public host: %s", req.URL.Hostname())
		}
		return nil
	}
}

// ValidatePublicHTTPSURL requires an HTTPS URL whose host resolves exclusively
// to publicly routable addresses.
func ValidatePublicHTTPSURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid HTTPS URL: %w", err)
	}
	if parsed.Scheme != httpsScheme || parsed.Hostname() == "" {
		return errors.New("URL must use HTTPS and include a host")
	}
	if isHostnameInternal(parsed.Hostname(), false) {
		return fmt.Errorf("URL host %q is not publicly routable", parsed.Hostname())
	}
	return nil
}

// GetPublicHTTPClient returns an HTTP client that can connect only to publicly
// routable addresses and validates each redirect and resolved dial address.
func GetPublicHTTPClient() *http.Client {
	return newOutboundHTTPClient(false)
}

// GetFederationHTTPClient returns the restricted outbound client used for
// federation. The test-only internal-federation override is preserved.
func GetFederationHTTPClient() *http.Client {
	return newOutboundHTTPClient(AllowInternalFederation())
}

// GetRetryableHTTPClient returns the federation HTTP client with retry logic
// for transient 502, 503, and 504 responses.
func GetRetryableHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 100 * time.Millisecond
	retryClient.RetryWaitMax = 1 * time.Second
	retryClient.Logger = nil
	retryClient.HTTPClient = GetFederationHTTPClient()
	client := retryClient.StandardClient()
	client.Timeout = 8 * time.Second
	return client
}
