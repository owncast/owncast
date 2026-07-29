package router

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

// Automatic HTTPS: when enabled (OWNCAST_ENABLE_AUTO_HTTPS +
// OWNCAST_HOST_NAME, currently set by packaged installs), Owncast attempts
// to open ports 80 and 443 next to the regular web server port and obtains
// a Let's Encrypt certificate on demand via ACME. The regular plain HTTP
// listener is unaffected in every case: any failure here is logged and the
// server continues HTTP-only.
//
// Certificates are issued lazily by autocert on the first TLS handshake for
// the configured hostname and cached in the data directory, so restarts do
// not contact the certificate authority. acmeFailureBackoff guards Let's
// Encrypt's failed-validation rate limit (~5/hour per hostname) against
// repeated handshakes (e.g. an operator refreshing their browser) while DNS
// or reachability is still misconfigured.
const acmeFailureBackoff = 20 * time.Minute

type autoTLS struct {
	host    string
	manager *autocert.Manager
	// lookup resolves a certificate for a handshake; it is
	// manager.GetCertificate outside of tests.
	lookup func(*tls.ClientHelloInfo) (*tls.Certificate, error)

	mu          sync.Mutex
	lastFailure time.Time
	certLogged  bool
}

// startAutoHTTPS attempts to bring up the automatic HTTPS listeners,
// serving handler over TLS on port 443 with an ACME challenge + redirect
// listener on port 80. It never fails the caller.
func startAutoHTTPS(cfg *config.Config, handler http.Handler) {
	httpPort := envPortOr("OWNCAST_TLS_HTTP_PORT", 80)
	httpsPort := envPortOr("OWNCAST_TLS_HTTPS_PORT", 443)

	httpListener, httpsListener, err := bindAutoTLSListeners(cfg.WebServerIP, httpPort, httpsPort)
	if err != nil {
		logAutoTLSBindFailure(err)
		return
	}

	at := newAutoTLS(cfg.AutoHTTPSHost)
	go at.logDNSDiagnostic()

	challengeServer := &http.Server{
		ReadHeaderTimeout: 4 * time.Second,
		Handler:           at.manager.HTTPHandler(nil),
	}
	go func() {
		if err := challengeServer.Serve(httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warnln("Automatic HTTPS: the port 80 challenge listener stopped:", err)
		}
	}()

	tlsConfig := at.manager.TLSConfig()
	tlsConfig.GetCertificate = at.getCertificate
	httpsServer := &http.Server{
		ReadHeaderTimeout: 4 * time.Second,
		Handler:           handler,
		TLSConfig:         tlsConfig,
	}
	go func() {
		if err := httpsServer.ServeTLS(httpsListener, "", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warnln("Automatic HTTPS: the port 443 listener stopped:", err)
		}
	}()

	log.Infof("Automatic HTTPS: listening for %s. A certificate will be requested with the first request.", at.host)
}

func newAutoTLS(host string) *autoTLS {
	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache(filepath.Join(config.DataDirectory, "certs")),
	}

	// Test and staging hook: point ACME at Pebble or the Let's Encrypt
	// staging environment instead of production.
	if directory := os.Getenv("OWNCAST_ACME_DIRECTORY"); directory != "" {
		client := &acme.Client{DirectoryURL: directory}
		if tlsConfig := utils.GetTLSConfig(); tlsConfig != nil {
			client.HTTPClient = &http.Client{
				Transport: &http.Transport{TLSClientConfig: tlsConfig},
			}
		}
		manager.Client = client
	}

	at := &autoTLS{host: host, manager: manager}
	at.lookup = manager.GetCertificate
	return at
}

// getCertificate wraps autocert's certificate lookup with an in-memory
// failure backoff. Cached certificates are served without any certificate
// authority contact; only failed issuance attempts for the configured
// hostname arm the backoff. Handshakes for other server names are rejected
// by the host policy and never affect it.
func (at *autoTLS) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if len(hello.SupportedProtos) == 1 && hello.SupportedProtos[0] == acme.ALPNProto {
		return at.lookup(hello)
	}

	isConfiguredHost := strings.EqualFold(strings.TrimSuffix(hello.ServerName, "."), at.host)

	if isConfiguredHost {
		at.mu.Lock()
		backingOff := !at.lastFailure.IsZero() && time.Since(at.lastFailure) < acmeFailureBackoff
		at.mu.Unlock()
		if backingOff {
			return nil, errors.New("automatic HTTPS: pausing certificate requests after a recent failure")
		}
	}

	cert, err := at.lookup(hello)

	if !isConfiguredHost {
		return cert, err
	}

	at.mu.Lock()
	defer at.mu.Unlock()
	if err != nil {
		at.lastFailure = time.Now()
		log.Errorf("Automatic HTTPS: getting a certificate for %s failed: %v. Check that DNS for this hostname points at this server and that port 80 is reachable from the internet. The next attempt will be made in %v.", at.host, err, acmeFailureBackoff)
		return nil, err
	}
	at.lastFailure = time.Time{}
	if !at.certLogged {
		at.certLogged = true
		log.Infof("Automatic HTTPS: certificate for %s is active.", at.host)
	}
	return cert, nil
}

// bindAutoTLSListeners opens the challenge and TLS listeners together.
// Binding is all-or-nothing: a half-working setup (challenges without TLS
// or vice versa) would be harder to reason about than no TLS at all.
func bindAutoTLSListeners(ip string, httpPort, httpsPort int) (net.Listener, net.Listener, error) {
	httpListener, err := net.Listen("tcp", net.JoinHostPort(ip, strconv.Itoa(httpPort)))
	if err != nil {
		return nil, nil, err
	}

	httpsListener, err := net.Listen("tcp", net.JoinHostPort(ip, strconv.Itoa(httpsPort)))
	if err != nil {
		_ = httpListener.Close()
		return nil, nil, err
	}

	return httpListener, httpsListener, nil
}

func logAutoTLSBindFailure(err error) {
	switch {
	case errors.Is(err, syscall.EADDRINUSE):
		log.Infoln("Automatic HTTPS: ports 80 and/or 443 are already in use, so automatic HTTPS is disabled. If you are running a reverse proxy this is expected.")
	case errors.Is(err, syscall.EACCES), errors.Is(err, syscall.EPERM):
		log.Warnln("Automatic HTTPS: insufficient permission to open ports 80 and 443, so automatic HTTPS is disabled. Binding low ports requires the CAP_NET_BIND_SERVICE capability or equivalent.")
	default:
		log.Warnln("Automatic HTTPS: unable to open ports 80 and 443, so automatic HTTPS is disabled:", err)
	}
}

// logDNSDiagnostic is an advisory-only check for the most common
// first-boot problem: the hostname not resolving yet. It never blocks
// certificate issuance — local resolution can fail (split-horizon DNS,
// NAT) in setups where external ACME validation would still succeed.
func (at *autoTLS) logDNSDiagnostic() {
	const recheckInterval = 10 * time.Minute
	for attempt := 0; ; attempt++ {
		if _, err := net.LookupHost(at.host); err == nil {
			if attempt > 0 {
				log.Infof("Automatic HTTPS: %s now resolves in DNS.", at.host)
			}
			return
		}
		log.Warnf("Automatic HTTPS: %s does not resolve in DNS yet. HTTPS will start working once DNS for this hostname points at this server. Checking again in %v.", at.host, recheckInterval)
		time.Sleep(recheckInterval)
	}
}

// envPortOr supports overriding the privileged default ports in tests.
func envPortOr(name string, defaultPort int) int {
	value := os.Getenv(name)
	if value == "" {
		return defaultPort
	}
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		log.Warnf("Ignoring invalid %s value %q.", name, value)
		return defaultPort
	}
	return port
}
