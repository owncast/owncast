package utils

import (
	"crypto/tls"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
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

// GetRetryableHTTPClient returns an http.Client with retry logic for transient failures.
// It uses hashicorp/go-retryablehttp with exponential backoff for 502, 503, 504 errors.
func GetRetryableHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 100 * time.Millisecond
	retryClient.RetryWaitMax = 1 * time.Second
	retryClient.Logger = nil // Disable default logging

	// Configure transport with connection pooling limits
	transport := GetHTTPTransportWithTLS(&http.Transport{
		MaxIdleConns:        20,               // Limit resource usage
		MaxIdleConnsPerHost: 2,                // Conservative per-host limit
		IdleConnTimeout:     10 * time.Second, // Fast cleanup of idle connections
		DisableKeepAlives:   false,
	})
	retryClient.HTTPClient.Transport = transport

	client := retryClient.StandardClient()
	client.Timeout = 8 * time.Second // Short timeout - legitimate servers respond quickly
	return client
}
