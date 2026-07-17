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

	if !d.allowInternal {
		public := publicIPAddresses(addresses)
		if len(public) == 0 {
			return nil, fmt.Errorf("refusing to connect to host %q: it resolves only to non-public addresses", host)
		}
		addresses = public
	}
	addresses = interleaveAddressFamilies(addresses)

	var dialErrors []error
	for i, address := range addresses {
		attemptCtx, cancel := attemptContext(ctx, len(addresses)-i)
		conn, err := d.dial(attemptCtx, network, net.JoinHostPort(address.IP.String(), port))
		cancel()
		if err == nil {
			return conn, nil
		}
		dialErrors = append(dialErrors, err)
	}

	return nil, fmt.Errorf("unable to connect to outbound host %q: %w", host, errors.Join(dialErrors...))
}

// attemptContext splits the remaining deadline evenly across the remaining
// dial attempts so one unreachable address cannot consume the entire
// request budget.
func attemptContext(ctx context.Context, attemptsLeft int) (context.Context, context.CancelFunc) {
	deadline, ok := ctx.Deadline()
	if !ok || attemptsLeft <= 1 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, time.Until(deadline)/time.Duration(attemptsLeft))
}

// interleaveAddressFamilies alternates address families, keeping the
// resolver's preferred family first, so a broken path for one family falls
// back to the other on the next attempt instead of after every address of
// the first family has timed out.
func interleaveAddressFamilies(addresses []net.IPAddr) []net.IPAddr {
	var v4, v6 []net.IPAddr
	for _, address := range addresses {
		if address.IP.To4() != nil {
			v4 = append(v4, address)
		} else {
			v6 = append(v6, address)
		}
	}
	if len(v4) == 0 || len(v6) == 0 {
		return addresses
	}
	first, second := v6, v4
	if addresses[0].IP.To4() != nil {
		first, second = v4, v6
	}
	interleaved := make([]net.IPAddr, 0, len(addresses))
	for i := 0; i < len(first) || i < len(second); i++ {
		if i < len(first) {
			interleaved = append(interleaved, first[i])
		}
		if i < len(second) {
			interleaved = append(interleaved, second[i])
		}
	}
	return interleaved
}

func newOutboundHTTPClient(allowInternal bool) *http.Client {
	transport := GetHTTPTransportWithTLS(&http.Transport{
		DialContext:         newRestrictedDialer(allowInternal).DialContext,
		ForceAttemptHTTP2:   true,
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
