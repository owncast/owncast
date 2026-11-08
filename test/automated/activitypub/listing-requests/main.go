//go:build integration

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-fed/httpsig"
)

const (
	adminUser = "admin"
	adminPass = "abc123"
	actorPath = "/federation/user/"
)

type activity struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Actor     json.RawMessage `json:"actor"`
	Object    json.RawMessage `json:"object"`
	Target    json.RawMessage `json:"target"`
	Directory bool            `json:"https://owncast.online/ns#directory"`
	To        json.RawMessage `json:"to"`
}

type actorDocument struct {
	ID        string `json:"id"`
	Inbox     string `json:"inbox"`
	PublicKey struct {
		ID           string `json:"id"`
		Owner        string `json:"owner"`
		PublicKeyPEM string `json:"publicKeyPem"`
	} `json:"publicKey"`
}

type serverRow struct {
	ID           int    `json:"id"`
	IRI          string `json:"iri"`
	FollowStatus string `json:"followStatus"`
	IsOnline     *bool  `json:"isOnline"`
}

type capturedRequest struct {
	Method    string
	URL       string
	Host      string
	Path      string
	UserAgent string
	Header    http.Header
	Body      []byte
	Activity  activity
	At        time.Time
}

type httpsProxy struct {
	server          *http.Server
	baseURL         string
	block           atomic.Bool
	blockedAttempts atomic.Int64
	mu              sync.RWMutex
	captures        []capturedRequest
	fakePath        string
	fakeBody        []byte
}

func newHTTPSProxy(cert tls.Certificate, backendPort int) (*httpsProxy, error) {
	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
	if err != nil {
		return nil, err
	}
	p := &httpsProxy{
		baseURL: "https://" + listener.Addr().String(),
	}
	reverse := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http", Host: net.JoinHostPort("127.0.0.1", strconv.Itoa(backendPort)),
	})
	director := reverse.Director
	reverse.Director = func(req *http.Request) {
		originalHost := req.Host
		director(req)
		req.Host = originalHost
		req.Header.Set("X-Forwarded-Proto", "https")
	}
	p.server = &http.Server{
		Handler:           http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { p.serve(reverse, w, r) }),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		if err := p.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("proxy %s stopped: %v", p.baseURL, err)
		}
	}()
	return p, nil
}

func (p *httpsProxy) serve(reverse http.Handler, w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	fakePath, fakeBody := p.fakePath, p.fakeBody
	p.mu.RUnlock()
	if r.Method == http.MethodGet && fakePath != "" && r.URL.Path == fakePath {
		w.Header().Set("Content-Type", "application/activity+json")
		_, _ = w.Write(fakeBody)
		return
	}

	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(io.LimitReader(r.Body, 8<<20))
		if err != nil {
			http.Error(w, "could not read proxied request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	if r.Method == http.MethodPost && len(body) > 0 {
		var message activity
		if json.Unmarshal(body, &message) == nil && message.Type != "" {
			captured := capturedRequest{
				Method:    r.Method,
				URL:       "https://" + r.Host + r.URL.RequestURI(),
				Host:      r.Host,
				Path:      r.URL.RequestURI(),
				UserAgent: r.UserAgent(),
				Header:    r.Header.Clone(),
				Body:      body,
				Activity:  message,
				At:        time.Now(),
			}
			p.mu.Lock()
			p.captures = append(p.captures, captured)
			p.mu.Unlock()
		}
	}
	if r.Method == http.MethodPost && fakePath != "" && r.URL.Path == fakePath+"/inbox" {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if p.block.Load() && r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/inbox") {
		p.blockedAttempts.Add(1)
		http.Error(w, "test receiver is temporarily unavailable", http.StatusServiceUnavailable)
		return
	}
	reverse.ServeHTTP(w, r)
}

func (p *httpsProxy) setFakeActor(actorPath string, body []byte) {
	p.mu.Lock()
	p.fakePath = actorPath
	p.fakeBody = append([]byte(nil), body...)
	p.mu.Unlock()
}

func (p *httpsProxy) snapshot() []capturedRequest {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]capturedRequest, len(p.captures))
	copy(result, p.captures)
	return result
}

func (p *httpsProxy) close(ctx context.Context) {
	_ = p.server.Shutdown(ctx)
}

type genericDirectory struct {
	server       *http.Server
	baseURL      string
	actorID      string
	inbox        string
	publicKeyPEM string
	mu           sync.RWMutex
	captures     []capturedRequest
}

func newGenericDirectory(cert tls.Certificate) (*genericDirectory, error) {
	// ponytail: this receiver never signs; its test TLS key also supplies actor metadata.
	publicDER, err := x509.MarshalPKIXPublicKey(&cert.PrivateKey.(*rsa.PrivateKey).PublicKey)
	if err != nil {
		return nil, err
	}
	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
	if err != nil {
		return nil, err
	}
	d := &genericDirectory{
		baseURL:      "https://" + listener.Addr().String(),
		publicKeyPEM: string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})),
	}
	d.actorID = d.baseURL + "/actors/listing-directory"
	d.inbox = d.baseURL + "/actors/listing-directory/inbox"
	d.server = &http.Server{Handler: http.HandlerFunc(d.serve), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := d.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("generic directory stopped: %v", err)
		}
	}()
	return d, nil
}

func (d *genericDirectory) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	switch r.URL.Path {
	case "/.well-known/nodeinfo":
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"links": []map[string]string{{
			"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
			"href": d.baseURL + "/nodeinfo/2.0",
		}}})
	case "/nodeinfo/2.0":
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"version":   "2.0",
			"software":  map[string]string{"name": "generic-listing-service", "version": "1"},
			"protocols": []string{"activitypub"},
			"metadata":  map[string]any{"federation": map[string]string{"username": "listing"}},
		})
	case "/.well-known/webfinger":
		w.Header().Set("Content-Type", "application/jrd+json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"subject": "acct:listing@" + r.Host,
			"links":   []map[string]string{{"rel": "self", "type": "application/activity+json", "href": d.actorID}},
		})
	case "/actors/listing-directory":
		w.Header().Set("Content-Type", "application/activity+json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"@context":          []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
			"id":                d.actorID,
			"type":              "Service",
			"preferredUsername": "listing",
			"name":              "Generic Listing Service",
			"inbox":             d.inbox,
			"publicKey": map[string]string{
				"id": d.actorID + "#main-key", "owner": d.actorID, "publicKeyPem": d.publicKeyPEM,
			},
		})
	case "/actors/listing-directory/inbox":
		body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
		if err != nil {
			http.Error(w, "could not read activity", http.StatusBadRequest)
			return
		}
		var message activity
		if json.Unmarshal(body, &message) != nil || message.Type == "" {
			http.Error(w, "invalid activity", http.StatusBadRequest)
			return
		}
		d.mu.Lock()
		d.captures = append(d.captures, capturedRequest{
			Method: r.Method, URL: "https://" + r.Host + r.URL.RequestURI(), Host: r.Host,
			Path: r.URL.RequestURI(), UserAgent: r.UserAgent(), Header: r.Header.Clone(),
			Body: body, Activity: message, At: time.Now(),
		})
		d.mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	default:
		http.NotFound(w, r)
	}
}

func (d *genericDirectory) snapshot() []capturedRequest {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]capturedRequest, len(d.captures))
	copy(result, d.captures)
	return result
}

func (d *genericDirectory) close(ctx context.Context) {
	_ = d.server.Shutdown(ctx)
}

type process struct {
	cmd     *exec.Cmd
	done    chan struct{}
	exitErr error
}

type instance struct {
	name      string
	bin       string
	workDir   string
	webPort   int
	rtmpPort  int
	proxy     *httpsProxy
	process   *process
	logPath   string
	username  string
	publicURL string
	database  string
}

func newInstance(name, bin, root, username string, cert tls.Certificate) (*instance, error) {
	webPort, err := freePort()
	if err != nil {
		return nil, err
	}
	rtmpPort, err := freePort()
	if err != nil {
		return nil, err
	}
	proxy, err := newHTTPSProxy(cert, webPort)
	if err != nil {
		return nil, err
	}
	workDir := filepath.Join(root, name)
	if err := os.MkdirAll(workDir, 0o700); err != nil {
		proxy.close(context.Background())
		return nil, err
	}
	return &instance{
		name: name, bin: bin, workDir: workDir, webPort: webPort, rtmpPort: rtmpPort,
		proxy: proxy, logPath: filepath.Join(root, name+".log"), username: username,
		publicURL: proxy.baseURL, database: filepath.Join(workDir, "owncast.db"),
	}, nil
}

func (i *instance) start() error {
	file, err := os.OpenFile(i.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	cmd := exec.Command(i.bin,
		"-database", i.database,
		"-webserverport", fmt.Sprint(i.webPort),
		"-webserverip", "127.0.0.1",
		"-rtmpport", fmt.Sprint(i.rtmpPort),
	)
	cmd.Dir = i.workDir
	cmd.Env = owncastEnv()
	cmd.Stdout = file
	cmd.Stderr = file
	if err := cmd.Start(); err != nil {
		_ = file.Close()
		return fmt.Errorf("start Owncast %s: %w", i.name, err)
	}
	proc := &process{cmd: cmd, done: make(chan struct{})}
	i.process = proc
	go func() {
		proc.exitErr = cmd.Wait()
		_ = file.Close()
		close(proc.done)
	}()
	return nil
}

func owncastEnv() []string {
	env := make([]string, 0, len(os.Environ())+2)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "OWNCAST_ALLOW_INTERNAL_FEDERATION=") || strings.HasPrefix(entry, "OWNCAST_INSECURE_SKIP_VERIFY=") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "OWNCAST_ALLOW_INTERNAL_FEDERATION=true", "OWNCAST_INSECURE_SKIP_VERIFY=true")
}

func (i *instance) stop() {
	if i.process == nil || i.process.cmd.Process == nil {
		return
	}
	proc := i.process
	if err := proc.cmd.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
		_ = proc.cmd.Process.Kill()
	}
	select {
	case <-proc.done:
	case <-time.After(10 * time.Second):
		_ = proc.cmd.Process.Kill()
		<-proc.done
	}
	i.process = nil
}

func (i *instance) apiURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", i.webPort)
}

func hostPort(rawURL string) string {
	parsed, _ := url.Parse(rawURL)
	return parsed.Host
}

func (i *instance) actorURL() string {
	return i.publicURL + actorPath + i.username
}

func (i *instance) request(ctx context.Context, method, endpoint string, body any) (int, []byte, error) {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, i.apiURL()+endpoint, reader)
	if err != nil {
		return 0, nil, err
	}
	req.SetBasicAuth(adminUser, adminPass)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := apiClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	return resp.StatusCode, response, err
}

func (i *instance) get(ctx context.Context, endpoint string, out any) error {
	status, body, err := i.request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("GET %s%s returned HTTP %d: %s", i.name, endpoint, status, body)
	}
	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("decode GET %s%s response %q: %w", i.name, endpoint, body, err)
		}
	}
	return nil
}

func (i *instance) post(ctx context.Context, endpoint string, body any) error {
	return i.simple(ctx, http.MethodPost, endpoint, body)
}

func requestFeaturedFollow(ctx context.Context, requester, target *instance) error {
	if err := requester.post(ctx, "/api/admin/federation/servers", map[string]any{"url": target.publicURL}); err != nil {
		return fmt.Errorf("%s request to feature %s: %w", requester.name, target.name, err)
	}
	return nil
}

func (i *instance) listingDecision(ctx context.Context, id int, approved, expectedSuccess bool) error {
	endpoint := fmt.Sprintf("/api/admin/federation/servers/%d/listing-request", id)
	status, response, err := i.request(ctx, http.MethodPost, endpoint, map[string]any{"approved": approved})
	if err != nil {
		return err
	}
	var result struct {
		Success *bool   `json:"success"`
		Message *string `json:"message"`
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("decode listing decision %s response %q: %w", endpoint, response, err)
	}
	if result.Success == nil || result.Message == nil {
		return fmt.Errorf("listing decision %s response is missing success/message: %q", endpoint, response)
	}
	if *result.Success != expectedSuccess {
		return fmt.Errorf("listing decision %s returned success=%t message=%q, expected %t (HTTP %d)", endpoint, *result.Success, *result.Message, expectedSuccess, status)
	}
	if status >= 500 {
		return fmt.Errorf("listing decision %s returned HTTP %d: %s", endpoint, status, response)
	}
	return nil
}

func (i *instance) simple(ctx context.Context, method, endpoint string, body any) error {
	status, response, err := i.request(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("%s %s%s returned HTTP %d: %s", method, i.name, endpoint, status, response)
	}
	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("decode %s %s%s response %q: %w", method, i.name, endpoint, response, err)
	}
	if !result.Success {
		return fmt.Errorf("%s %s%s failed: %s", method, i.name, endpoint, result.Message)
	}
	return nil
}

var (
	apiClient        = &http.Client{Timeout: 20 * time.Second}
	federationClient = &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Test-only TLS endpoints use a generated local certificate.
		}},
	}
)

func freePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func makeTestCertificate() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}
	now := time.Now()
	template := x509.Certificate{
		SerialNumber:          big.NewInt(now.UnixNano()),
		NotBefore:             now.Add(-time.Minute),
		NotAfter:              now.Add(24 * time.Hour),
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return tls.X509KeyPair(certPEM, keyPEM)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	root := os.Getenv("OWNCAST_LISTING_TEST_ARTIFACTS")
	if root == "" {
		var err error
		root, err = os.MkdirTemp("", "owncast-listing-requests-")
		if err != nil {
			log.Fatalf("create test artifact directory: %v", err)
		}
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		log.Fatalf("create artifact directory %s: %v", root, err)
	}
	bin := os.Getenv("OWNCAST_LISTING_TEST_BIN")
	if bin == "" {
		log.Fatalf("OWNCAST_LISTING_TEST_BIN is not set; run test-listing-requests.sh")
	}
	if err := run(ctx, root, bin); err != nil {
		log.Printf("FAIL: %v", err)
		log.Printf("Artifacts preserved at %s", root)
		for _, name := range []string{"a.log", "b.log", "ffmpeg.log"} {
			printLog(filepath.Join(root, name))
		}
		os.Exit(1)
	}
	log.Printf("PASS: listing-request end-to-end scenarios completed; artifacts: %s", root)
}

func printLog(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil || len(data) == 0 {
		return
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) > 120 {
		lines = lines[len(lines)-120:]
	}
	log.Printf("--- %s (last %d lines) ---\n%s", filename, len(lines), strings.Join(lines, "\n"))
}

func run(ctx context.Context, root, bin string) (result error) {
	cert, err := makeTestCertificate()
	if err != nil {
		return fmt.Errorf("create local TLS certificate: %w", err)
	}
	a, err := newInstance("a", bin, root, "streamera", cert)
	if err != nil {
		return err
	}
	b, err := newInstance("b", bin, root, "streamerb", cert)
	if err != nil {
		a.proxy.close(context.Background())
		return err
	}
	defer func() {
		a.stop()
		b.stop()
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.proxy.close(closeCtx)
		b.proxy.close(closeCtx)
	}()
	directory, err := newGenericDirectory(cert)
	if err != nil {
		return fmt.Errorf("start generic HTTPS directory: %w", err)
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		directory.close(closeCtx)
	}()

	log.Printf("Starting isolated Owncast A and B processes (scratch: %s)", root)
	if err := a.start(); err != nil {
		return err
	}
	if err := b.start(); err != nil {
		return err
	}
	if err := waitReady(ctx, a); err != nil {
		return err
	}
	if err := waitReady(ctx, b); err != nil {
		return err
	}
	if err := configure(ctx, a); err != nil {
		return err
	}
	if err := configure(ctx, b); err != nil {
		return err
	}
	if err := assertDefaultAutoAcceptFalse(ctx, b); err != nil {
		return err
	}
	if err := assertDefaultAutoAcceptFalse(ctx, a); err != nil {
		return err
	}
	if err := checkHTTPSAdminWrite(ctx, a); err != nil {
		return err
	}
	aActor, err := getActor(ctx, a.actorURL())
	if err != nil {
		return fmt.Errorf("fetch A actor: %w", err)
	}
	bActor, err := getActor(ctx, b.actorURL())
	if err != nil {
		return fmt.Errorf("fetch B actor: %w", err)
	}
	if aActor.ID != a.actorURL() || bActor.ID != b.actorURL() {
		return fmt.Errorf("unexpected actor IDs: A=%q B=%q", aActor.ID, bActor.ID)
	}

	fansA, err := fanSnapshot(ctx, a)
	if err != nil {
		return fmt.Errorf("snapshot A fans: %w", err)
	}
	fansB, err := fanSnapshot(ctx, b)
	if err != nil {
		return fmt.Errorf("snapshot B fans: %w", err)
	}

	if err := negativeIngressChecks(ctx, a, b, aActor, bActor); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "negative ingress controls"); err != nil {
		return err
	}
	if err := negativeListingTargetChecks(ctx, a, b); err != nil {
		return err
	}

	log.Printf("Scenario: hostname-only request, exact consent survives A restart, concurrent B approval, receiver retry")
	b.proxy.block.Store(true)
	if err := a.post(ctx, "/api/admin/federation/listing-requests", map[string]any{"url": hostPort(b.publicURL)}); err != nil {
		return fmt.Errorf("create hostname-only listing request from A: %w", err)
	}
	var originalInvite capturedRequest
	if err := waitFor(ctx, "first signed Invite attempt to B", 15*time.Second, func() (bool, string) {
		found := findCapture(b.proxy.snapshot(), func(c capturedRequest) bool {
			return c.Activity.Type == "Invite" && idOf(c.Activity.Actor) == aActor.ID && isOwncastRequest(c)
		})
		if found == nil {
			return false, "no captured Invite yet"
		}
		originalInvite = *found
		return true, "captured"
	}); err != nil {
		return err
	}
	if err := assertInvite(originalInvite.Activity, aActor.ID, bActor.ID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(originalInvite, aActor); err != nil {
		return fmt.Errorf("verify producer-signed Invite from A: %w", err)
	}
	if b.proxy.blockedAttempts.Load() == 0 {
		return fmt.Errorf("receiver-outage scenario did not force an inbox delivery failure")
	}
	// B's inbox has returned a real 503. Restart B on the same scratch DB, then
	// restore delivery so A's durable outbox can retry the captured activity.
	b.stop()
	if err := b.start(); err != nil {
		return fmt.Errorf("restart receiver B: %w", err)
	}
	if err := waitReady(ctx, b); err != nil {
		return err
	}
	b.proxy.block.Store(false)
	if err := waitForServerStatus(ctx, b, a.publicURL, "requested", 120*time.Second); err != nil {
		return fmt.Errorf("wait for queued Invite after B restart: %w", err)
	}
	if err := assertRequestedVisibleOnlyInAdminServers(ctx, b, a.publicURL); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "incoming Invite"); err != nil {
		return err
	}

	if err := replayCaptured(ctx, originalInvite); err != nil {
		return fmt.Errorf("replay captured duplicate Invite: %w", err)
	}
	if err := waitFor(ctx, "duplicate Invite handling", 5*time.Second, func() (bool, string) {
		rows, err := adminServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		if countIRI(rows, a.publicURL) != 1 {
			return false, fmt.Sprintf("B has %d rows for A", countIRI(rows, a.publicURL))
		}
		return true, "one row"
	}); err != nil {
		return err
	}
	bRow, err := findServerRow(ctx, b, a.publicURL)
	if err != nil {
		return err
	}
	log.Printf("B listing request row id=%d status=%s", bRow.ID, bRow.FollowStatus)
	a.stop()
	if err := a.start(); err != nil {
		return fmt.Errorf("restart requester A with pending listing consent: %w", err)
	}
	if err := waitReady(ctx, a); err != nil {
		return err
	}
	if err := assertAutoAcceptConfig(ctx, a, false); err != nil {
		return fmt.Errorf("A's unrelated global auto-accept setting changed across restart: %w", err)
	}
	forgedActorIRI, forgedFollowID, err := sendSameHostDirectoryFollow(ctx, b, aActor)
	if err != nil {
		return err
	}
	if err := waitFor(ctx, "same-host forged Follow remains pending", 30*time.Second, func() (bool, string) {
		if err := assertFeatureRequest(ctx, a, forgedActorIRI, true); err != nil {
			return false, err.Error()
		}
		return true, forgedActorIRI
	}); err != nil {
		return err
	}
	if err := assertFeatureRequest(ctx, a, bActor.ID, false); err != nil {
		return err
	}
	if err := assertDirectoryFollower(ctx, a, forgedActorIRI, false); err != nil {
		return err
	}
	if got := activityCount(b.proxy.snapshot(), "Accept", aActor.ID, forgedFollowID); got != 0 {
		return fmt.Errorf("same-host forged Follow received %d Accept deliveries before the exact target followed", got)
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "same-host forged Follow"); err != nil {
		return err
	}
	if err := decideListingConcurrently(ctx, b, bRow.ID, true); err != nil {
		return err
	}
	if err := waitForServerStatus(ctx, b, a.publicURL, "accepted", 40*time.Second); err != nil {
		return err
	}
	inviteID := originalInvite.Activity.ID
	acceptInvite, followDirectory, err := waitApprovalActivities(ctx, a, bActor.ID, aActor.ID, inviteID, 40*time.Second)
	if err != nil {
		return err
	}
	if err := assertAccept(acceptInvite.Activity, bActor.ID, aActor.ID, inviteID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(acceptInvite, bActor); err != nil {
		return fmt.Errorf("verify producer-signed Accept from B: %w", err)
	}
	if err := assertFollow(followDirectory.Activity, bActor.ID, aActor.ID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(followDirectory, bActor); err != nil {
		return fmt.Errorf("verify producer-signed directory Follow from B: %w", err)
	}
	followAccept, err := waitForActivity(ctx, b, "Accept", aActor.ID, followDirectory.Activity.ID, 40*time.Second)
	if err != nil {
		return fmt.Errorf("A did not automatically accept B's exact marked Follow: %w", err)
	}
	if err := assertFollowAccept(followAccept.Activity, aActor.ID, bActor.ID, followDirectory.Activity.ID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(followAccept, aActor); err != nil {
		return fmt.Errorf("verify producer-signed automatic Accept from A: %w", err)
	}
	if err := assertCountForActivity(a.proxy.snapshot(), "Accept", bActor.ID, inviteID, 1); err != nil {
		return err
	}
	if err := assertCountForActivity(a.proxy.snapshot(), "Follow", bActor.ID, aActor.ID, 1); err != nil {
		return err
	}
	if got := activityCount(b.proxy.snapshot(), "Accept", aActor.ID, forgedFollowID); got != 0 {
		return fmt.Errorf("same-host forged Follow received %d Accept deliveries after B's exact actor followed", got)
	}
	if err := assertCountForActivity(b.proxy.snapshot(), "Accept", aActor.ID, followDirectory.Activity.ID, 1); err != nil {
		return err
	}
	if err := b.listingDecision(ctx, bRow.ID, false, false); err != nil {
		return fmt.Errorf("opposite decision must return success=false without changing approval: %w", err)
	}
	if err := waitForServerStatus(ctx, b, a.publicURL, "accepted", 5*time.Second); err != nil {
		return err
	}
	if err := assertFeatureRequest(ctx, a, bActor.ID, false); err != nil {
		return err
	}
	if err := assertFeatureRequest(ctx, a, forgedActorIRI, true); err != nil {
		return err
	}
	if err := assertDirectoryFollower(ctx, a, bActor.ID, true); err != nil {
		return err
	}
	if err := assertDirectoryFollower(ctx, a, forgedActorIRI, false); err != nil {
		return err
	}
	if err := assertSinglePublicListing(ctx, b, a.publicURL); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "manual B approval completed exact consent"); err != nil {
		return err
	}

	log.Printf("Scenario: remove listing, verify replay/consumed consent cannot revive relationship")
	if err := removeRelationship(ctx, a, b, bActor.ID); err != nil {
		return err
	}
	beforeAccepts := activityCount(a.proxy.snapshot(), "Accept", bActor.ID, inviteID)
	beforeFollows := activityCount(a.proxy.snapshot(), "Follow", bActor.ID, aActor.ID)
	if beforeAccepts != 1 || beforeFollows != 1 {
		return fmt.Errorf("before replay expected one Accept and Follow, got Accept=%d Follow=%d", beforeAccepts, beforeFollows)
	}
	if err := replayCaptured(ctx, originalInvite); err != nil {
		return fmt.Errorf("replay consumed Invite after removal: %w", err)
	}
	if err := assertRemains(ctx, "consumed Invite remains removed", 3*time.Second, func() (bool, string) {
		rows, err := adminServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		if countIRI(rows, a.publicURL) != 0 {
			return false, "replayed Invite recreated B's listing row"
		}
		return true, "no row"
	}); err != nil {
		return err
	}
	if got := activityCount(a.proxy.snapshot(), "Accept", bActor.ID, inviteID); got != beforeAccepts {
		return fmt.Errorf("replayed consumed Invite caused %d Accept deliveries; before=%d", got, beforeAccepts)
	}
	if got := activityCount(a.proxy.snapshot(), "Follow", bActor.ID, aActor.ID); got != beforeFollows {
		return fmt.Errorf("replayed consumed Invite caused %d Follow deliveries; before=%d", got, beforeFollows)
	}
	log.Printf("Scenario: fresh marked Follow without remaining consent stays pending")
	followStart := len(a.proxy.snapshot())
	if err := requestFeaturedFollow(ctx, b, a); err != nil {
		return err
	}
	var unsolicitedFollow capturedRequest
	if err := waitFor(ctx, "fresh B marked Follow reaches A", 40*time.Second, func() (bool, string) {
		captures := a.proxy.snapshot()
		for index := followStart; index < len(captures); index++ {
			captured := captures[index]
			if captured.Activity.Type == "Follow" && idOf(captured.Activity.Actor) == bActor.ID && objectID(captured.Activity) == aActor.ID {
				unsolicitedFollow = captured
				return true, "captured"
			}
		}
		return false, "no fresh Follow yet"
	}); err != nil {
		return err
	}
	if err := assertFollow(unsolicitedFollow.Activity, bActor.ID, aActor.ID); err != nil {
		return err
	}
	if err := assertRemains(ctx, "consumed consent does not auto-accept a fresh Follow", 4*time.Second, func() (bool, string) {
		if err := assertFeatureRequest(ctx, a, bActor.ID, true); err != nil {
			return false, err.Error()
		}
		if err := assertDirectoryFollower(ctx, a, bActor.ID, false); err != nil {
			return false, err.Error()
		}
		if got := activityCount(b.proxy.snapshot(), "Accept", aActor.ID, unsolicitedFollow.Activity.ID); got != 0 {
			return false, fmt.Sprintf("fresh Follow received %d Accept deliveries", got)
		}
		return true, "pending and not accepted"
	}); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "fresh unsolicited marked Follow"); err != nil {
		return err
	}
	if err := removeRelationship(ctx, a, b, bActor.ID); err != nil {
		return err
	}
	if err := waitFor(ctx, "fresh pending Follow removed", 20*time.Second, func() (bool, string) {
		if err := assertFeatureRequest(ctx, a, bActor.ID, false); err != nil {
			return false, err.Error()
		}
		return true, "no pending feature request"
	}); err != nil {
		return err
	}

	log.Printf("Scenario: B auto-accepts Invite and A automatically accepts the exact marked Follow")
	if err := setAutoAccept(ctx, b, true); err != nil {
		return err
	}
	if err := a.post(ctx, "/api/admin/federation/listing-requests", map[string]any{"url": hostPort(b.publicURL)}); err != nil {
		return fmt.Errorf("create hostname-only auto-accept listing request from A: %w", err)
	}
	if err := waitForServerStatus(ctx, b, a.publicURL, "accepted", 60*time.Second); err != nil {
		return fmt.Errorf("auto-accepted listing request did not reach accepted: %w", err)
	}
	autoInvite, err := findLatestInvite(b.proxy.snapshot(), aActor.ID)
	if err != nil {
		return err
	}
	if err := assertInvite(autoInvite.Activity, aActor.ID, bActor.ID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(autoInvite, aActor); err != nil {
		return fmt.Errorf("verify producer-signed auto-accepted Invite: %w", err)
	}
	autoAccept, autoFollow, err := waitApprovalActivities(ctx, a, bActor.ID, aActor.ID, autoInvite.Activity.ID, 40*time.Second)
	if err != nil {
		return err
	}
	if err := assertAccept(autoAccept.Activity, bActor.ID, aActor.ID, autoInvite.Activity.ID); err != nil {
		return err
	}
	if err := assertFollow(autoFollow.Activity, bActor.ID, aActor.ID); err != nil {
		return err
	}
	if err := waitForServerStatus(ctx, b, a.publicURL, "accepted", 40*time.Second); err != nil {
		return err
	}
	autoFollowAccept, err := waitForActivity(ctx, b, "Accept", aActor.ID, autoFollow.Activity.ID, 40*time.Second)
	if err != nil {
		return fmt.Errorf("A did not automatically accept B's marked Follow after auto-approval: %w", err)
	}
	if err := assertFollowAccept(autoFollowAccept.Activity, aActor.ID, bActor.ID, autoFollow.Activity.ID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(autoFollowAccept, aActor); err != nil {
		return fmt.Errorf("verify producer-signed auto-scenario Accept from A: %w", err)
	}
	if err := assertFeatureRequest(ctx, a, bActor.ID, false); err != nil {
		return err
	}
	if err := assertDirectoryFollower(ctx, a, bActor.ID, true); err != nil {
		return err
	}
	if err := assertSinglePublicListing(ctx, b, a.publicURL); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "automatic B Invite approval completed exact consent"); err != nil {
		return err
	}

	log.Printf("Scenario: real A stream is live, then offline, as observed in B's public listing")
	if err := verifyStreamLifecycle(ctx, root, a, b, a.publicURL); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "stream lifecycle"); err != nil {
		return err
	}

	log.Printf("Scenario: B rejects a later Invite")
	if err := removeRelationship(ctx, a, b, bActor.ID); err != nil {
		return err
	}
	if err := setAutoAccept(ctx, b, false); err != nil {
		return err
	}
	followsBeforeReject := activityCount(a.proxy.snapshot(), "Follow", bActor.ID, aActor.ID)
	if err := a.post(ctx, "/api/admin/federation/listing-requests", map[string]any{"url": hostPort(b.publicURL)}); err != nil {
		return fmt.Errorf("create hostname-only request for rejection scenario: %w", err)
	}
	if err := waitForServerStatus(ctx, b, a.publicURL, "requested", 40*time.Second); err != nil {
		return err
	}
	rejectedRow, err := findServerRow(ctx, b, a.publicURL)
	if err != nil {
		return err
	}
	rejectedInvite, err := findLatestInvite(b.proxy.snapshot(), aActor.ID)
	if err != nil {
		return err
	}
	if err := b.post(ctx, fmt.Sprintf("/api/admin/federation/servers/%d/listing-request", rejectedRow.ID), map[string]any{"approved": false}); err != nil {
		return fmt.Errorf("B rejects listing request: %w", err)
	}
	if err := waitForServerStatus(ctx, b, a.publicURL, "rejected", 40*time.Second); err != nil {
		return err
	}
	rejectActivity, err := waitForActivity(ctx, a, "Reject", bActor.ID, rejectedInvite.Activity.ID, 40*time.Second)
	if err != nil {
		return err
	}
	if err := assertReject(rejectActivity.Activity, bActor.ID, aActor.ID, rejectedInvite.Activity.ID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(rejectActivity, bActor); err != nil {
		return fmt.Errorf("verify producer-signed Reject from B: %w", err)
	}
	if err := assertCountForActivity(a.proxy.snapshot(), "Reject", bActor.ID, rejectedInvite.Activity.ID, 1); err != nil {
		return err
	}
	if got := activityCount(a.proxy.snapshot(), "Accept", bActor.ID, rejectedInvite.Activity.ID); got != 0 {
		return fmt.Errorf("rejected Invite %s received %d Accept deliveries", rejectedInvite.Activity.ID, got)
	}
	if got := activityCount(a.proxy.snapshot(), "Follow", bActor.ID, aActor.ID); got != followsBeforeReject {
		return fmt.Errorf("rejecting Invite caused another directory Follow: before=%d after=%d", followsBeforeReject, got)
	}
	if err := assertNotPubliclyListed(ctx, b, a.publicURL); err != nil {
		return err
	}
	if err := assertFeatureRequest(ctx, a, bActor.ID, false); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "rejected listing request"); err != nil {
		return err
	}
	if err := removeRelationship(ctx, a, b, bActor.ID); err != nil {
		return err
	}
	rejectedFollowStart := len(a.proxy.snapshot())
	if err := requestFeaturedFollow(ctx, b, a); err != nil {
		return fmt.Errorf("request a fresh B-to-A listing after Invite rejection: %w", err)
	}
	var postRejectFollow capturedRequest
	if err := waitFor(ctx, "fresh B marked Follow after Invite rejection", 40*time.Second, func() (bool, string) {
		captures := a.proxy.snapshot()
		for index := rejectedFollowStart; index < len(captures); index++ {
			captured := captures[index]
			if captured.Activity.Type == "Follow" && idOf(captured.Activity.Actor) == bActor.ID && objectID(captured.Activity) == aActor.ID {
				postRejectFollow = captured
				return true, "captured"
			}
		}
		return false, "no fresh Follow yet"
	}); err != nil {
		return err
	}
	if err := assertFollow(postRejectFollow.Activity, bActor.ID, aActor.ID); err != nil {
		return err
	}
	if err := assertRemains(ctx, "Reject(Invite) revoked the outgoing consent", 4*time.Second, func() (bool, string) {
		if err := assertFeatureRequest(ctx, a, bActor.ID, true); err != nil {
			return false, err.Error()
		}
		if err := assertDirectoryFollower(ctx, a, bActor.ID, false); err != nil {
			return false, err.Error()
		}
		if got := activityCount(b.proxy.snapshot(), "Accept", aActor.ID, postRejectFollow.Activity.ID); got != 0 {
			return false, fmt.Sprintf("post-rejection Follow received %d Accept deliveries", got)
		}
		return true, "fresh Follow remains pending"
	}); err != nil {
		return err
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "fresh Follow after Invite rejection"); err != nil {
		return err
	}
	if err := removeRelationship(ctx, a, b, bActor.ID); err != nil {
		return err
	}
	if err := waitFor(ctx, "post-rejection Follow removed", 20*time.Second, func() (bool, string) {
		if err := assertFeatureRequest(ctx, a, bActor.ID, false); err != nil {
			return false, err.Error()
		}
		return true, "no pending feature request"
	}); err != nil {
		return err
	}

	log.Printf("Scenario: hostname-only generic NodeInfo/WebFinger discovery without Owncast capability")
	genericStart := len(directory.snapshot())
	if err := a.post(ctx, "/api/admin/federation/listing-requests", map[string]any{"url": hostPort(directory.baseURL)}); err != nil {
		return fmt.Errorf("create listing request from generic NodeInfo/WebFinger host: %w", err)
	}
	genericInvite, err := waitForGenericInviteSince(ctx, directory, genericStart, 60*time.Second)
	if err != nil {
		return err
	}
	if err := assertInvite(genericInvite.Activity, aActor.ID, directory.actorID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(genericInvite, aActor); err != nil {
		return fmt.Errorf("verify producer-signed generic-directory Invite: %w", err)
	}
	if err := assertFanSnapshots(ctx, a, b, fansA, fansB, "generic directory request"); err != nil {
		return err
	}

	log.Printf("PASS: hostname discovery, manual/automatic Invite consent, exact actor binding, restart persistence, rejection revocation, consumed-consent removal, malformed inputs, no fan effects, typeless Invite, stream live/offline, and receiver retry")
	if os.Getenv("KEEP_RUNNING") == "true" {
		log.Printf("KEEP_RUNNING: frontend=current UI embedded in both binaries; A_ADMIN=%s/admin/config-featured/ B_ADMIN=%s/admin/config-featured/ A_HOSTNAME=%s B_HOSTNAME=%s ADMIN_USER=%s ADMIN_PASS=%s", a.apiURL(), b.apiURL(), hostPort(a.publicURL), hostPort(b.publicURL), adminUser, adminPass)
		<-ctx.Done()
	}
	return nil
}

func configure(ctx context.Context, i *instance) error {
	log.Printf("Configuring %s at %s", i.name, i.publicURL)
	settings := []struct {
		endpoint string
		value    any
	}{
		{"/api/admin/config/serverurl", i.publicURL},
		{"/api/admin/config/federation/username", i.username},
		{"/api/admin/config/federation/private", false},
		{"/api/admin/config/federation/enable", true},
	}
	for _, setting := range settings {
		if err := i.post(ctx, setting.endpoint, map[string]any{"value": setting.value}); err != nil {
			return fmt.Errorf("configure %s: %w", setting.endpoint, err)
		}
	}
	return nil
}

func waitReady(ctx context.Context, i *instance) error {
	return waitFor(ctx, "Owncast "+i.name+" HTTP readiness", 60*time.Second, func() (bool, string) {
		requestCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, i.apiURL()+"/api/status", nil)
		if err != nil {
			return false, err.Error()
		}
		resp, err := apiClient.Do(req)
		if err != nil {
			if i.process != nil {
				select {
				case <-i.process.done:
					return false, fmt.Sprintf("process exited: %v", i.process.exitErr)
				default:
				}
			}
			return false, err.Error()
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return false, fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		return true, "ready"
	})
}

func waitFor(ctx context.Context, label string, timeout time.Duration, check func() (bool, string)) error {
	deadline := time.Now().Add(timeout)
	last := "not checked"
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("%s interrupted: %w", label, err)
		}
		ok, detail := check()
		if ok {
			log.Printf("PASS: %s (%s)", label, detail)
			return nil
		}
		last = detail
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s interrupted: %w", label, ctx.Err())
		case <-ticker.C:
		}
	}
	return fmt.Errorf("timed out after %s waiting for %s; last state: %s", timeout, label, last)
}

func assertRemains(ctx context.Context, label string, duration time.Duration, check func() (bool, string)) error {
	deadline := time.NewTimer(duration)
	defer deadline.Stop()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	for {
		if ok, detail := check(); !ok {
			return fmt.Errorf("%s: %s", label, detail)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			log.Printf("PASS: %s for %s", label, duration)
			return nil
		case <-ticker.C:
		}
	}
}

func assertDefaultAutoAcceptFalse(ctx context.Context, i *instance) error {
	var config map[string]any
	if err := i.get(ctx, "/api/admin/serverconfig", &config); err != nil {
		return fmt.Errorf("read fresh %s admin server config: %w", i.name, err)
	}
	federation, ok := config["federation"].(map[string]any)
	if !ok {
		return fmt.Errorf("%s server config lacks federation object: %s", i.name, mustJSON(config))
	}
	value, ok := federation["autoAcceptListingRequests"].(bool)
	if !ok {
		return fmt.Errorf("%s server config lacks boolean federation.autoAcceptListingRequests: %s", i.name, mustJSON(config))
	}
	if value {
		return fmt.Errorf("fresh %s instance unexpectedly auto-accepts listing requests", i.name)
	}
	log.Printf("PASS: fresh %s autoAcceptListingRequests defaults false", i.name)
	return nil
}

func setAutoAccept(ctx context.Context, i *instance, enabled bool) error {
	if err := i.post(ctx, "/api/admin/config/federation/autoacceptlistingrequests", map[string]any{"value": enabled}); err != nil {
		return fmt.Errorf("set %s autoAcceptListingRequests=%t: %w", i.name, enabled, err)
	}
	return assertAutoAcceptConfig(ctx, i, enabled)
}

func checkHTTPSAdminWrite(ctx context.Context, i *instance) error {
	if err := setAutoAcceptHTTPS(ctx, i, true); err != nil {
		return fmt.Errorf("HTTPS Origin+Basic admin write: %w", err)
	}
	if err := setAutoAcceptHTTPS(ctx, i, false); err != nil {
		return fmt.Errorf("restore autoAcceptListingRequests=false through HTTPS: %w", err)
	}
	log.Printf("PASS: HTTPS same-origin Basic-auth admin write changed server config and restored autoAcceptListingRequests=false")
	return nil
}

func setAutoAcceptHTTPS(ctx context.Context, i *instance, enabled bool) error {
	body, err := json.Marshal(map[string]any{"value": enabled})
	if err != nil {
		return err
	}
	endpoint := "/api/admin/config/federation/autoacceptlistingrequests"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.publicURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(adminUser, adminPass)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", i.publicURL)
	resp, err := federationClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("decode HTTPS admin response %q: %w", response, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !result.Success {
		return fmt.Errorf("HTTPS POST %s returned HTTP %d success=%t: %s", endpoint, resp.StatusCode, result.Success, result.Message)
	}
	return assertAutoAcceptConfig(ctx, i, enabled)
}

func assertAutoAcceptConfig(ctx context.Context, i *instance, enabled bool) error {
	var config map[string]any
	if err := i.get(ctx, "/api/admin/serverconfig", &config); err != nil {
		return err
	}
	federation, _ := config["federation"].(map[string]any)
	actual, ok := federation["autoAcceptListingRequests"].(bool)
	if !ok || actual != enabled {
		return fmt.Errorf("%s autoAcceptListingRequests is %t after setting %t (present=%t)", i.name, actual, enabled, ok)
	}
	return nil
}

func getActor(ctx context.Context, actorURL string) (actorDocument, error) {
	var actor actorDocument
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, actorURL, nil)
	if err != nil {
		return actor, err
	}
	req.Header.Set("Accept", "application/activity+json")
	resp, err := federationClient.Do(req)
	if err != nil {
		return actor, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return actor, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return actor, fmt.Errorf("GET %s returned HTTP %d: %s", actorURL, resp.StatusCode, body)
	}
	if err := json.Unmarshal(body, &actor); err != nil {
		return actor, fmt.Errorf("decode actor %s %q: %w", actorURL, body, err)
	}
	if actor.ID == "" || actor.Inbox == "" {
		return actor, fmt.Errorf("actor %s lacks id or inbox: %s", actorURL, body)
	}
	return actor, nil
}

func adminServers(ctx context.Context, i *instance) ([]serverRow, error) {
	var response struct {
		Servers []serverRow `json:"servers"`
	}
	if err := i.get(ctx, "/api/admin/federation/servers", &response); err != nil {
		return nil, err
	}
	return response.Servers, nil
}

func publicServers(ctx context.Context, i *instance) ([]serverRow, error) {
	var response struct {
		Servers []serverRow `json:"servers"`
	}
	if err := i.get(ctx, "/api/federation/servers", &response); err != nil {
		return nil, err
	}
	return response.Servers, nil
}

func countIRI(rows []serverRow, iri string) int {
	count := 0
	for _, row := range rows {
		if row.IRI == iri {
			count++
		}
	}
	return count
}

func findServerRow(ctx context.Context, i *instance, iri string) (serverRow, error) {
	rows, err := adminServers(ctx, i)
	if err != nil {
		return serverRow{}, fmt.Errorf("read %s admin servers: %w", i.name, err)
	}
	for _, row := range rows {
		if row.IRI == iri {
			return row, nil
		}
	}
	return serverRow{}, fmt.Errorf("%s admin servers has no row for %s", i.name, iri)
}

func waitForServerStatus(ctx context.Context, i *instance, iri, status string, timeout time.Duration) error {
	return waitFor(ctx, fmt.Sprintf("%s row %s has followStatus=%s", i.name, iri, status), timeout, func() (bool, string) {
		rows, err := adminServers(ctx, i)
		if err != nil {
			return false, err.Error()
		}
		count := countIRI(rows, iri)
		for _, row := range rows {
			if row.IRI == iri && row.FollowStatus == status && count == 1 {
				return true, fmt.Sprintf("row id=%d status=%s", row.ID, row.FollowStatus)
			}
		}
		return false, fmt.Sprintf("matching rows=%d requested status=%q", count, status)
	})
}

func assertRequestedVisibleOnlyInAdminServers(ctx context.Context, i *instance, iri string) error {
	rows, err := adminServers(ctx, i)
	if err != nil {
		return err
	}
	if countIRI(rows, iri) != 1 {
		return fmt.Errorf("incoming listing request should have exactly one admin server row, got %d", countIRI(rows, iri))
	}
	row, err := findServerRow(ctx, i, iri)
	if err != nil {
		return err
	}
	if row.FollowStatus != "requested" {
		return fmt.Errorf("new incoming listing request has status %q, expected requested", row.FollowStatus)
	}
	var featureRequests struct {
		Requests []struct {
			Link string `json:"link"`
		} `json:"requests"`
	}
	if err := i.get(ctx, "/api/admin/federation/feature-requests", &featureRequests); err != nil {
		return err
	}
	for _, request := range featureRequests.Requests {
		if strings.HasPrefix(request.Link, iri+"/") {
			return fmt.Errorf("listing request %s leaked into the separate feature-requests endpoint", iri)
		}
	}
	public, err := publicServers(ctx, i)
	if err != nil {
		return err
	}
	if countIRI(public, iri) != 0 {
		return fmt.Errorf("unapproved listing request %s leaked into public federation servers", iri)
	}
	return nil
}

func assertFeatureRequest(ctx context.Context, i *instance, actorIRI string, expected bool) error {
	var response struct {
		Requests []struct {
			Link string `json:"link"`
		} `json:"requests"`
	}
	if err := i.get(ctx, "/api/admin/federation/feature-requests", &response); err != nil {
		return err
	}
	found := false
	for _, request := range response.Requests {
		found = found || request.Link == actorIRI
	}
	if found != expected {
		return fmt.Errorf("%s feature-requests contains actor %s=%t; expected %t", i.name, actorIRI, found, expected)
	}
	return nil
}

func assertNotPubliclyListed(ctx context.Context, i *instance, iri string) error {
	rows, err := publicServers(ctx, i)
	if err != nil {
		return err
	}
	if countIRI(rows, iri) != 0 {
		return fmt.Errorf("%s public federation list includes %s before manual approval", i.name, iri)
	}
	return nil
}

func assertSinglePublicListing(ctx context.Context, i *instance, iri string) error {
	rows, err := publicServers(ctx, i)
	if err != nil {
		return err
	}
	if countIRI(rows, iri) != 1 {
		return fmt.Errorf("%s public federation list has %d entries for %s, expected exactly one; response=%s", i.name, countIRI(rows, iri), iri, mustJSON(rows))
	}
	row, err := findServerRow(ctx, i, iri)
	if err != nil {
		return err
	}
	if row.FollowStatus != "accepted" {
		return fmt.Errorf("%s admin row for %s has followStatus=%q, expected accepted", i.name, iri, row.FollowStatus)
	}
	return nil
}

func fanSnapshot(ctx context.Context, i *instance) (json.RawMessage, error) {
	var response json.RawMessage
	if err := i.get(ctx, "/api/admin/followers?limit=1000&offset=0", &response); err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(response, &decoded); err != nil {
		return nil, err
	}
	canonical, err := json.Marshal(decoded)
	return canonical, err
}

func assertFanSnapshots(ctx context.Context, a, b *instance, expectedA, expectedB json.RawMessage, stage string) error {
	actualA, err := fanSnapshot(ctx, a)
	if err != nil {
		return fmt.Errorf("read A fan list/count after %s: %w", stage, err)
	}
	if !bytes.Equal(actualA, expectedA) {
		return fmt.Errorf("A fan list/count changed during %s: before=%s after=%s", stage, expectedA, actualA)
	}
	actualB, err := fanSnapshot(ctx, b)
	if err != nil {
		return fmt.Errorf("read B fan list/count after %s: %w", stage, err)
	}
	if !bytes.Equal(actualB, expectedB) {
		return fmt.Errorf("B fan list/count changed during %s: before=%s after=%s", stage, expectedB, actualB)
	}
	return nil
}

func assertDirectoryFollower(ctx context.Context, i *instance, actorIRI string, expected bool) error {
	var followers []map[string]any
	if err := i.get(ctx, "/api/admin/followers/directory", &followers); err != nil {
		return err
	}
	count := 0
	for _, follower := range followers {
		for _, key := range []string{"link", "iri", "actorIRI"} {
			if actor, _ := follower[key].(string); actor == actorIRI {
				count++
				break
			}
		}
	}
	if (count == 1) != expected {
		return fmt.Errorf("%s directory follower list has %d rows for %s, expected present=%t: %s", i.name, count, actorIRI, expected, mustJSON(followers))
	}
	return nil
}

func decideListingConcurrently(ctx context.Context, i *instance, id int, approved bool) error {
	endpoint := fmt.Sprintf("/api/admin/federation/servers/%d/listing-request", id)
	const workers = 5
	results := make(chan error, workers)
	var group sync.WaitGroup
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			results <- i.post(ctx, endpoint, map[string]any{"approved": approved})
		}()
	}
	group.Wait()
	close(results)
	for err := range results {
		if err != nil {
			return fmt.Errorf("same listing decision must be idempotent under concurrency: %w", err)
		}
	}
	// Sequential retries also catch handlers that only happen to be safe under a race.
	for range 2 {
		if err := i.post(ctx, endpoint, map[string]any{"approved": approved}); err != nil {
			return fmt.Errorf("same listing decision was not idempotent on retry: %w", err)
		}
	}
	return nil
}

func waitApprovalActivities(ctx context.Context, receiver *instance, actorIRI, objectIRI, inviteID string, timeout time.Duration) (capturedRequest, capturedRequest, error) {
	var accept, follow capturedRequest
	err := waitFor(ctx, "Accept(Invite) and directory Follow deliveries", timeout, func() (bool, string) {
		for _, captured := range receiver.proxy.snapshot() {
			if !isOwncastRequest(captured) {
				continue
			}
			if captured.Activity.Type == "Accept" && idOf(captured.Activity.Actor) == actorIRI && objectID(captured.Activity) == inviteID {
				accept = captured
			}
			if captured.Activity.Type == "Follow" && idOf(captured.Activity.Actor) == actorIRI && objectID(captured.Activity) == objectIRI {
				follow = captured
			}
		}
		if accept.Activity.ID != "" && follow.Activity.ID != "" {
			return true, "captured signed Accept and Follow"
		}
		return false, fmt.Sprintf("Accept id=%q Follow id=%q", accept.Activity.ID, follow.Activity.ID)
	})
	return accept, follow, err
}

func removeRelationship(ctx context.Context, a, b *instance, bActor string) error {
	row, err := findServerRow(ctx, b, a.publicURL)
	if err != nil {
		return fmt.Errorf("find B server row before removal: %w", err)
	}
	if err := b.simple(ctx, http.MethodDelete, fmt.Sprintf("/api/admin/federation/servers/%d", row.ID), nil); err != nil {
		return fmt.Errorf("remove B server listing row: %w", err)
	}
	if err := waitFor(ctx, "relationship removal", 20*time.Second, func() (bool, string) {
		rows, err := adminServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		if countIRI(rows, a.publicURL) != 0 {
			return false, fmt.Sprintf("B retains %d listing rows", countIRI(rows, a.publicURL))
		}
		if err := assertDirectoryFollower(ctx, a, bActor, false); err != nil {
			return false, err.Error()
		}
		return true, "B listing row and A directory follower removed"
	}); err != nil {
		return err
	}
	return nil
}

// assertUnsolicitedAcceptCannotPublish verifies that a valid signed Accept(Follow)
// for a Follow that was never sent cannot promote a requested or rejected listing.
func assertUnsolicitedAcceptCannotPublish(ctx context.Context, b *instance, fakeActorID string, fakeKey *rsa.PrivateKey, bActor actorDocument, fixtureRow serverRow) error {
	for _, status := range []string{"requested", "rejected"} {
		if status == "rejected" {
			if err := b.post(ctx, fmt.Sprintf("/api/admin/federation/servers/%d/listing-request", fixtureRow.ID), map[string]any{"approved": false}); err != nil {
				return fmt.Errorf("reject typeless Invite to reach rejected state: %w", err)
			}
			if err := waitForServerStatus(ctx, b, fixtureRow.IRI, "rejected", 20*time.Second); err != nil {
				return err
			}
		}
		accept := map[string]any{
			"@context": "https://www.w3.org/ns/activitystreams",
			"id":       fakeActorID + "/accept/" + randomID(),
			"type":     "Accept",
			"actor":    fakeActorID,
			"object": map[string]any{
				"id":    b.publicURL + "/federation/follow/never-sent",
				"type":  "Follow",
				"actor": bActor.ID,
			},
			"to":                                     []string{bActor.ID},
			"https://owncast.online/ns#streamStatus": "live",
		}
		if err := sendSignedActivity(ctx, bActor.Inbox, accept, fakeKey, fakeActorID+"#main-key", false); err != nil {
			return fmt.Errorf("send unsolicited Accept while %s: %w", status, err)
		}
		if err := assertRemains(ctx, "unsolicited Accept cannot override "+status+" listing", 3*time.Second, func() (bool, string) {
			row, err := findServerRow(ctx, b, fixtureRow.IRI)
			if err != nil {
				return false, err.Error()
			}
			if row.FollowStatus != status || (row.IsOnline != nil && *row.IsOnline) {
				return false, fmt.Sprintf("unsolicited Accept changed listing: %s", mustJSON(row))
			}
			if err := assertNotPubliclyListed(ctx, b, fixtureRow.IRI); err != nil {
				return false, err.Error()
			}
			return true, "listing remains " + status + " and unpublished"
		}); err != nil {
			return err
		}
	}
	log.Printf("PASS: signed unsolicited Accept(Follow) cannot publish requested or rejected listings")
	return nil
}

func negativeIngressChecks(ctx context.Context, a, b *instance, aActor, bActor actorDocument) error {
	fakeKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate negative-control key: %w", err)
	}
	fakeActorID := a.publicURL + "/test/fake-actor"
	publicDER, err := x509.MarshalPKIXPublicKey(&fakeKey.PublicKey)
	if err != nil {
		return fmt.Errorf("encode negative-control public key: %w", err)
	}
	fakePEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER}))
	fakeActor, err := json.Marshal(map[string]any{
		"@context":          []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
		"id":                fakeActorID,
		"type":              "Service",
		"preferredUsername": "streamera",
		"inbox":             fakeActorID + "/inbox",
		"publicKey": map[string]string{
			"id": fakeActorID + "#main-key", "owner": fakeActorID, "publicKeyPem": fakePEM,
		},
	})
	if err != nil {
		return err
	}
	a.proxy.setFakeActor("/test/fake-actor", fakeActor)

	beforeA := countOwncastActivities(a.proxy.snapshot())
	beforeB := countOwncastActivities(b.proxy.snapshot())
	invalidInvite := makeInvite(aActor.ID, aActor.ID, bActor.ID, bActor.ID)
	if err := sendSignedActivity(ctx, bActor.Inbox, invalidInvite, fakeKey, aActor.PublicKey.ID, true); err != nil {
		return fmt.Errorf("send invalid-signature Invite: %w", err)
	}
	wrongTarget := makeInvite(fakeActorID, fakeActorID, aActor.ID, bActor.ID)
	if err := sendSignedActivity(ctx, bActor.Inbox, wrongTarget, fakeKey, fakeActorID+"#main-key", false); err != nil {
		return fmt.Errorf("send wrong-target Invite: %w", err)
	}
	impersonation := makeInvite(aActor.ID, aActor.ID, bActor.ID, bActor.ID)
	if err := sendSignedActivity(ctx, bActor.Inbox, impersonation, fakeKey, fakeActorID+"#main-key", false); err != nil {
		return fmt.Errorf("send same-host actor-impersonation Invite: %w", err)
	}
	if err := assertRemains(ctx, "negative Invites produce no relationship or federation delivery", 3*time.Second, func() (bool, string) {
		rows, err := adminServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		if countIRI(rows, a.publicURL) != 0 {
			return false, fmt.Sprintf("unexpected rows: %s", mustJSON(rows))
		}
		if countOwncastActivities(a.proxy.snapshot()) != beforeA || countOwncastActivities(b.proxy.snapshot()) != beforeB {
			return false, "Owncast emitted an AP response to a negative control"
		}
		return true, "no rows and no outgoing AP messages"
	}); err != nil {
		return fmt.Errorf("negative signature/target/actor-owner ingress checks failed: %w", err)
	}
	log.Printf("PASS: invalid signature, wrong target, and same-host/different-owner actor claim were rejected")
	rowsBeforeFixture, err := adminServers(ctx, b)
	if err != nil {
		return err
	}
	knownRows := make(map[int]struct{}, len(rowsBeforeFixture))
	for _, row := range rowsBeforeFixture {
		knownRows[row.ID] = struct{}{}
	}
	typelessInviteID := "urn:uuid:" + randomID()
	typelessInvite := map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       typelessInviteID,
		"type":     "Invite",
		"actor":    map[string]string{"id": fakeActorID},
		"object":   map[string]string{"id": fakeActorID},
		"target":   map[string]string{"id": bActor.ID},
		"to":       []string{bActor.ID},
	}
	if err := sendSignedActivity(ctx, bActor.Inbox, typelessInvite, fakeKey, fakeActorID+"#main-key", false); err != nil {
		return fmt.Errorf("send valid typeless-field Invite: %w", err)
	}
	var fixtureRow serverRow
	if err := waitFor(ctx, "valid typeless Invite creates one requested admin server row", 30*time.Second, func() (bool, string) {
		rows, err := adminServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		var added []serverRow
		for _, row := range rows {
			if _, exists := knownRows[row.ID]; !exists {
				added = append(added, row)
			}
		}
		if len(added) != 1 {
			return false, fmt.Sprintf("new rows=%s", mustJSON(added))
		}
		fixtureRow = added[0]
		return fixtureRow.FollowStatus == "requested", fmt.Sprintf("new row id=%d iri=%s status=%s", fixtureRow.ID, fixtureRow.IRI, fixtureRow.FollowStatus)
	}); err != nil {
		return fmt.Errorf("typeless embedded actor/object/target IDs were not accepted: %w", err)
	}
	if err := assertUnsolicitedAcceptCannotPublish(ctx, b, fakeActorID, fakeKey, bActor, fixtureRow); err != nil {
		return err
	}
	var fixtureReject capturedRequest
	if err := waitFor(ctx, "signed Reject delivered to typeless fixture actor", 30*time.Second, func() (bool, string) {
		for _, captured := range a.proxy.snapshot() {
			if isOwncastRequest(captured) && captured.Path == "/test/fake-actor/inbox" && captured.Activity.Type == "Reject" && idOf(captured.Activity.Actor) == bActor.ID && objectID(captured.Activity) == typelessInviteID {
				fixtureReject = captured
				return true, "captured"
			}
		}
		return false, "no Reject yet"
	}); err != nil {
		return err
	}
	if err := assertReject(fixtureReject.Activity, bActor.ID, fakeActorID, typelessInviteID); err != nil {
		return err
	}
	if err := verifyCapturedSignature(fixtureReject, bActor); err != nil {
		return fmt.Errorf("verify producer-signed Reject for typeless Invite: %w", err)
	}
	if got := activityCount(b.proxy.snapshot(), "Accept", bActor.ID, typelessInviteID); got != 0 {
		return fmt.Errorf("rejected typeless Invite received %d Accept deliveries", got)
	}
	if err := assertNotPubliclyListed(ctx, b, fixtureRow.IRI); err != nil {
		return err
	}
	if err := b.simple(ctx, http.MethodDelete, fmt.Sprintf("/api/admin/federation/servers/%d", fixtureRow.ID), nil); err != nil {
		return fmt.Errorf("remove typeless fixture row before real listing request: %w", err)
	}
	if err := waitFor(ctx, "typeless fixture row removed", 20*time.Second, func() (bool, string) {
		rows, err := adminServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		for _, row := range rows {
			if row.ID == fixtureRow.ID {
				return false, fmt.Sprintf("fixture row remains: %s", mustJSON(row))
			}
		}
		return true, "fixture row removed"
	}); err != nil {
		return err
	}
	log.Printf("PASS: a signed Invite with typeless embedded actor/object/target IDs was accepted and rejected")
	return nil
}

func negativeListingTargetChecks(ctx context.Context, requester, destination *instance) error {
	destinationHost := hostPort(destination.publicURL)
	inputs := []string{
		destination.publicURL,
		destinationHost + "/federation/user/" + destination.username,
		destinationHost + "?query=1",
		destinationHost + "#fragment",
		"user@" + destinationHost,
		destination.actorURL(),
	}
	beforeA := countOwncastActivities(requester.proxy.snapshot())
	beforeB := countOwncastActivities(destination.proxy.snapshot())
	for _, input := range inputs {
		status, body, err := requester.request(ctx, http.MethodPost, "/api/admin/federation/listing-requests", map[string]any{"url": input})
		if err != nil {
			return fmt.Errorf("submit invalid listing target %q: %w", input, err)
		}
		var response struct {
			Success bool `json:"success"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("decode invalid listing target %q response %q: %w", input, body, err)
		}
		if status != http.StatusBadRequest || response.Success {
			return fmt.Errorf("invalid listing target %q returned HTTP %d success=%t; expected HTTP 400 failure", input, status, response.Success)
		}
	}
	if countOwncastActivities(requester.proxy.snapshot()) != beforeA || countOwncastActivities(destination.proxy.snapshot()) != beforeB {
		return fmt.Errorf("invalid listing target input caused federation activity")
	}
	log.Printf("PASS: listing requests reject schemes, paths, queries, fragments, userinfo, and actor URLs")
	return nil
}

func sendSameHostDirectoryFollow(ctx context.Context, target *instance, receiverActor actorDocument) (string, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("generate same-host actor key: %w", err)
	}
	actorIRI := target.publicURL + "/test/forged-directory"
	publicDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("encode same-host actor key: %w", err)
	}
	actor, err := json.Marshal(map[string]any{
		"@context":          []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
		"id":                actorIRI,
		"type":              "Service",
		"preferredUsername": "forged-directory",
		"inbox":             actorIRI + "/inbox",
		"publicKey": map[string]string{
			"id": actorIRI + "#main-key", "owner": actorIRI,
			"publicKeyPem": string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})),
		},
	})
	if err != nil {
		return "", "", err
	}
	target.proxy.setFakeActor("/test/forged-directory", actor)
	followID := "urn:uuid:" + randomID()
	follow := map[string]any{
		"@context":                            "https://www.w3.org/ns/activitystreams",
		"id":                                  followID,
		"type":                                "Follow",
		"actor":                               actorIRI,
		"object":                              receiverActor.ID,
		"to":                                  []string{receiverActor.ID},
		"https://owncast.online/ns#directory": true,
	}
	if err := sendSignedActivity(ctx, receiverActor.Inbox, follow, key, actorIRI+"#main-key", false); err != nil {
		return "", "", fmt.Errorf("deliver same-host forged marked Follow: %w", err)
	}
	return actorIRI, followID, nil
}

func makeInvite(actorIRI, objectIRI, targetIRI, recipientIRI string) map[string]any {
	return map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       "urn:uuid:" + randomID(),
		"type":     "Invite",
		"actor":    actorIRI,
		"object":   objectIRI,
		"target":   targetIRI,
		"to":       []string{recipientIRI},
	}
}

func randomID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])
}

func sendSignedActivity(ctx context.Context, inbox string, message any, privateKey *rsa.PrivateKey, keyID string, corruptSignature bool) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, inbox, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/activity+json")
	req.Header.Set("Accept", `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Set("Host", req.URL.Host)
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.RSA_SHA256}, httpsig.DigestSha256,
		[]string{httpsig.RequestTarget, "host", "date", "digest"}, httpsig.Signature, 0,
	)
	if err != nil {
		return err
	}
	if err := signer.SignRequest(privateKey, keyID, req, body); err != nil {
		return err
	}
	if corruptSignature {
		req.Header.Set("Signature", `keyId="`+keyID+`",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="bm90LWEtdmFsaWQtc2lnbmF0dXJl"`)
	}
	req.Header.Set("User-Agent", "listing-request-e2e-negative-control")
	resp, err := federationClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("negative control received HTTP %d: %s", resp.StatusCode, response)
	}
	return nil
}

func replayCaptured(ctx context.Context, captured capturedRequest) error {
	req, err := http.NewRequestWithContext(ctx, captured.Method, captured.URL, bytes.NewReader(captured.Body))
	if err != nil {
		return err
	}
	req.Host = captured.Host
	for _, key := range []string{"Accept", "Content-Type", "Date", "Digest", "Signature", "User-Agent"} {
		for _, value := range captured.Header.Values(key) {
			req.Header.Add(key, value)
		}
	}
	resp, err := federationClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("captured replay received HTTP %d; relationship state is checked separately: %s", resp.StatusCode, response)
	}
	return nil
}

func findCapture(captures []capturedRequest, predicate func(capturedRequest) bool) *capturedRequest {
	for idx := range captures {
		if predicate(captures[idx]) {
			copy := captures[idx]
			return &copy
		}
	}
	return nil
}

func findLatestInvite(captures []capturedRequest, actorIRI string) (capturedRequest, error) {
	var result *capturedRequest
	for idx := range captures {
		captured := captures[idx]
		if captured.Activity.Type == "Invite" && idOf(captured.Activity.Actor) == actorIRI && isOwncastRequest(captured) {
			copy := captured
			if result == nil || copy.At.After(result.At) {
				result = &copy
			}
		}
	}
	if result == nil {
		return capturedRequest{}, fmt.Errorf("no producer-generated Invite by %s was captured", actorIRI)
	}
	return *result, nil
}

func waitForActivity(ctx context.Context, receiver *instance, kind, actorIRI, objectIRI string, timeout time.Duration) (capturedRequest, error) {
	var result capturedRequest
	err := waitFor(ctx, fmt.Sprintf("signed %s from %s for %s", kind, actorIRI, objectIRI), timeout, func() (bool, string) {
		for _, captured := range receiver.proxy.snapshot() {
			if isOwncastRequest(captured) && captured.Activity.Type == kind && idOf(captured.Activity.Actor) == actorIRI && objectID(captured.Activity) == objectIRI {
				result = captured
				return true, "captured"
			}
		}
		return false, "no matching activity yet"
	})
	return result, err
}

func waitForGenericInviteSince(ctx context.Context, d *genericDirectory, start int, timeout time.Duration) (capturedRequest, error) {
	var result capturedRequest
	err := waitFor(ctx, "Invite delivered to generic directory inbox", timeout, func() (bool, string) {
		captures := d.snapshot()
		for index := start; index < len(captures); index++ {
			if captures[index].Activity.Type == "Invite" {
				result = captures[index]
				return true, "captured"
			}
		}
		return false, "no Invite yet"
	})
	return result, err
}

func isOwncastRequest(c capturedRequest) bool {
	return strings.Contains(strings.ToLower(c.UserAgent), "owncast.online")
}

func countOwncastActivities(captures []capturedRequest) int {
	count := 0
	for _, captured := range captures {
		if isOwncastRequest(captured) {
			count++
		}
	}
	return count
}

func objectID(message activity) string {
	return idOf(message.Object)
}

func idOf(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var value string
	if json.Unmarshal(raw, &value) == nil {
		return value
	}
	var object struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(raw, &object) == nil {
		return object.ID
	}
	return ""
}

func containsIRI(raw json.RawMessage, expected string) bool {
	if idOf(raw) == expected {
		return true
	}
	var values []json.RawMessage
	if json.Unmarshal(raw, &values) == nil {
		for _, value := range values {
			if idOf(value) == expected {
				return true
			}
		}
	}
	return false
}

func assertInvite(message activity, actorIRI, targetIRI string) error {
	if message.Type != "Invite" || message.ID == "" || idOf(message.Actor) != actorIRI || idOf(message.Object) != actorIRI || idOf(message.Target) != targetIRI || !containsIRI(message.To, targetIRI) {
		return fmt.Errorf("Invite wire contract mismatch: got type=%q id=%q actor=%q object=%q target=%q to=%s", message.Type, message.ID, idOf(message.Actor), idOf(message.Object), idOf(message.Target), message.To)
	}
	return nil
}

func assertAccept(message activity, actorIRI, toIRI, inviteID string) error {
	if message.Type != "Accept" || message.ID == "" || idOf(message.Actor) != actorIRI || idOf(message.Object) != inviteID || !containsIRI(message.To, toIRI) {
		return fmt.Errorf("Accept(Invite) wire contract mismatch: type=%q id=%q actor=%q object=%q to=%s", message.Type, message.ID, idOf(message.Actor), idOf(message.Object), message.To)
	}
	return nil
}

func assertFollowAccept(message activity, actorIRI, followerIRI, followID string) error {
	if message.Type != "Accept" || message.ID == "" || idOf(message.Actor) != actorIRI {
		return fmt.Errorf("Accept(Follow) wire contract mismatch: type=%q id=%q actor=%q", message.Type, message.ID, idOf(message.Actor))
	}
	var follow activity
	if err := json.Unmarshal(message.Object, &follow); err != nil {
		return fmt.Errorf("decode accepted Follow: %w", err)
	}
	if follow.ID != followID {
		return fmt.Errorf("accepted Follow ID = %q, want %q", follow.ID, followID)
	}
	return assertFollow(follow, followerIRI, actorIRI)
}

func assertReject(message activity, actorIRI, toIRI, inviteID string) error {
	if message.Type != "Reject" || message.ID == "" || idOf(message.Actor) != actorIRI || idOf(message.Object) != inviteID || !containsIRI(message.To, toIRI) {
		return fmt.Errorf("Reject(Invite) wire contract mismatch: type=%q id=%q actor=%q object=%q to=%s", message.Type, message.ID, idOf(message.Actor), idOf(message.Object), message.To)
	}
	return nil
}

func assertFollow(message activity, actorIRI, directoryIRI string) error {
	if message.Type != "Follow" || message.ID == "" || idOf(message.Actor) != actorIRI || idOf(message.Object) != directoryIRI || !containsIRI(message.To, directoryIRI) || !message.Directory {
		return fmt.Errorf("marked directory Follow wire contract mismatch: type=%q id=%q actor=%q object=%q directory=%t to=%s", message.Type, message.ID, idOf(message.Actor), idOf(message.Object), message.Directory, message.To)
	}
	return nil
}

func activityCount(captures []capturedRequest, kind, actorIRI, objectIRI string) int {
	count := 0
	for _, captured := range captures {
		if isOwncastRequest(captured) && captured.Activity.Type == kind && idOf(captured.Activity.Actor) == actorIRI && objectID(captured.Activity) == objectIRI {
			count++
		}
	}
	return count
}

func assertCountForActivity(captures []capturedRequest, kind, actorIRI, objectIRI string, expected int) error {
	count := activityCount(captures, kind, actorIRI, objectIRI)
	if count != expected {
		return fmt.Errorf("captured %s deliveries from %s for %s: got %d, expected %d", kind, actorIRI, objectIRI, count, expected)
	}
	return nil
}

func verifyCapturedSignature(captured capturedRequest, actor actorDocument) error {
	block, _ := pem.Decode([]byte(actor.PublicKey.PublicKeyPEM))
	if block == nil {
		return fmt.Errorf("actor %s has invalid publicKeyPem", actor.ID)
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("actor %s public key is %T, not RSA", actor.ID, parsed)
	}
	req, err := http.NewRequest(captured.Method, captured.URL, bytes.NewReader(captured.Body))
	if err != nil {
		return err
	}
	req.Host = captured.Host
	req.Header = captured.Header.Clone()
	verifier, err := httpsig.NewVerifier(req)
	if err != nil {
		return fmt.Errorf("parse captured HTTP signature: %w", err)
	}
	if err := verifier.Verify(publicKey, httpsig.RSA_SHA256); err != nil {
		return fmt.Errorf("signature did not verify with %s's public key: %w", actor.ID, err)
	}
	return nil
}

func verifyStreamLifecycle(ctx context.Context, root string, a, b *instance, aIRI string) error {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is required for the live/offline scenario (install ffmpeg and retry): %w", err)
	}
	logPath := filepath.Join(root, "ffmpeg.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	cmd := exec.Command(ffmpeg,
		"-hide_banner", "-loglevel", "warning", "-re",
		"-f", "lavfi", "-i", "testsrc2=size=640x360:rate=15",
		"-f", "lavfi", "-i", "sine=frequency=1000:sample_rate=44100",
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-pix_fmt", "yuv420p", "-g", "30", "-c:a", "aac", "-b:a", "128k",
		"-f", "flv", fmt.Sprintf("rtmp://127.0.0.1:%d/live/abc123", a.rtmpPort),
	)
	cmd.Stdout, cmd.Stderr = logFile, logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return fmt.Errorf("start ffmpeg test stream: %w", err)
	}
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		_ = logFile.Close()
		close(done)
	}()
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Signal(os.Interrupt)
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				_ = cmd.Process.Kill()
				<-done
			}
		}
	}()
	if err := waitFor(ctx, "A source reports online", 60*time.Second, func() (bool, string) {
		var status map[string]any
		if err := a.get(ctx, "/api/status", &status); err != nil {
			return false, err.Error()
		}
		return status["online"] == true, fmt.Sprintf("A online=%v", status["online"])
	}); err != nil {
		return fmt.Errorf("stream source failed to become live: %w", err)
	}
	if err := waitFor(ctx, "B observes A live in public listing", 90*time.Second, func() (bool, string) {
		rows, err := publicServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		for _, row := range rows {
			if row.IRI == aIRI && row.IsOnline != nil && *row.IsOnline {
				return true, "B public row reports isOnline=true"
			}
		}
		return false, "B has not observed A online"
	}); err != nil {
		return err
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("stop ffmpeg test stream: %w", err)
	}
	select {
	case <-done:
	case <-time.After(8 * time.Second):
		_ = cmd.Process.Kill()
		<-done
	}
	if err := waitFor(ctx, "A source reports offline", 60*time.Second, func() (bool, string) {
		var status map[string]any
		if err := a.get(ctx, "/api/status", &status); err != nil {
			return false, err.Error()
		}
		return status["online"] == false, fmt.Sprintf("A online=%v", status["online"])
	}); err != nil {
		return fmt.Errorf("stream source did not become offline: %w", err)
	}
	if err := waitFor(ctx, "B observes A offline in public listing", 90*time.Second, func() (bool, string) {
		rows, err := publicServers(ctx, b)
		if err != nil {
			return false, err.Error()
		}
		for _, row := range rows {
			if row.IRI == aIRI && row.IsOnline != nil && !*row.IsOnline {
				return true, "B public row reports isOnline=false"
			}
		}
		return false, "B has not observed A offline"
	}); err != nil {
		return err
	}
	return nil
}

func mustJSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("<json error: %v>", err)
	}
	return string(encoded)
}
