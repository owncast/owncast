package utils

import (
	"errors"
	"net"
	"net/netip"
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

// CanonicalizeHostname returns a DNS-compatible ASCII hostname. It expects an
// unwrapped hostname, such as "::1", not URL host syntax, such as "[::1]".
func CanonicalizeHostname(hostname string) (string, error) {
	if hostname == "" {
		return "", errors.New("hostname is required")
	}

	if addr, err := netip.ParseAddr(hostname); err == nil {
		return strings.ToLower(addr.String()), nil
	}

	asciiHostname, err := idna.Lookup.ToASCII(hostname)
	if err != nil {
		return "", err
	}

	return strings.ToLower(asciiHostname), nil
}

// CanonicalizeHost returns a DNS-compatible ASCII host, preserving any port.
func CanonicalizeHost(host string) (string, error) {
	if strings.Contains(host, "://") {
		return "", errors.New("host appears to be a full URL; pass only the host[:port] portion")
	}

	u, err := url.Parse("//" + host)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		return "", errors.New("host is required")
	}
	if u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("host must not include user info, path, query, or fragment")
	}

	hostname, err := CanonicalizeHostname(u.Hostname())
	if err != nil {
		return "", err
	}

	return joinHostPort(hostname, u.Port()), nil
}

// CanonicalizeURLHostname returns a URL string with its hostname converted to
// DNS-compatible ASCII form.
func CanonicalizeURLHostname(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	canonicalURL, err := CanonicalizeURL(*u)
	if err != nil {
		return "", err
	}

	return canonicalURL.String(), nil
}

// CanonicalizeURL returns a copy of the provided URL with its hostname
// converted to DNS-compatible ASCII form.
func CanonicalizeURL(u url.URL) (*url.URL, error) {
	hostname, err := CanonicalizeHostname(u.Hostname())
	if err != nil {
		return nil, err
	}

	u.Host = joinHostPort(hostname, u.Port())

	return &u, nil
}

func joinHostPort(hostname string, port string) string {
	if port != "" {
		return net.JoinHostPort(hostname, port)
	}

	if strings.Contains(hostname, ":") {
		return "[" + hostname + "]"
	}

	return hostname
}
