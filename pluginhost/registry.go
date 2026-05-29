package pluginhost

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/services/plugins"
)

// Registry browse / install. By default the host points at the
// public catalog (DefaultPluginRegistry); operators can override
// with OWNCAST_PLUGIN_REGISTRY (set it to a private/self-hosted
// directory's API base, e.g. http://localhost:8088 for a local
// directory backend during development). The host appends bare
// resource paths to that base (`/plugins`, `/plugins/<slug>`), so
// the env var must point at whatever URL the registry's API answers
// from; the platform's /api/* rewrite or its absence is the
// operator's concern, not the host's.

// registryHTTPClient has a tight timeout because a slow or unreachable
// registry shouldn't hang the admin UI. The full install flow
// (list-fetch + .ocpkg download) needs to complete within a few
// requests' worth of time.
var registryHTTPClient = &http.Client{Timeout: 30 * time.Second}

// DefaultPluginRegistry is the URL the Owncast host points at when
// OWNCAST_PLUGIN_REGISTRY isn't set. It's the public catalog the
// Owncast project runs, so every Owncast instance gets a working
// Browse tab out of the box without per-deployment configuration.
// Operators who want a different catalog (a private/self-hosted
// directory, a staging instance, etc.) set the env var to override.
const DefaultPluginRegistry = "https://owncast.directory/api"

// registryBase reads the operator's configured registry URL from
// OWNCAST_PLUGIN_REGISTRY, falling back to the public catalog when
// unset. Trims any trailing slash so callers can safely concatenate
// path segments.
func registryBase() string {
	v := os.Getenv("OWNCAST_PLUGIN_REGISTRY")
	if v == "" {
		v = DefaultPluginRegistry
	}
	return strings.TrimRight(v, "/")
}

// handleRegistryRoute dispatches /api/admin/plugin-registry/<action>
// to the matching handler. Strips an optional trailing slash so the
// admin frontend's slash-canonicalizing redirect doesn't 404 against
// our routes. Registered as a prefix-handler in AdminHandler.
func (p *Host) handleRegistryRoute(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/admin/plugin-registry/"), "/")
	switch rest {
	case "list":
		p.handleRegistryList(w, r)
	case "install":
		p.handleRegistryInstall(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleRegistryList proxies the registry's public plugin listing.
// Used by the admin's "Browse plugins" view so the frontend talks
// only to its own origin (no extra CORS plumbing for the registry).
func (p *Host) handleRegistryList(w http.ResponseWriter, _ *http.Request) {
	base := registryBase()
	if base == "" {
		writeJSONResponse(w, http.StatusOK, []any{})
		return
	}
	resp, err := registryHTTPClient.Get(base + "/plugins")
	if err != nil {
		writeJSONResponse(w, http.StatusBadGateway, map[string]string{jsonErrorKey: "registry fetch: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 5<<20)) // 5 MiB cap on registry response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(body) //nolint:gosec // G705 false positive: bytes are JSON from a trusted upstream
}

// registryInstallRequest is the body shape posted by the admin UI.
// Slug identifies the plugin in the registry (the registry's primary
// key); version pins a specific release.
type registryInstallRequest struct {
	Slug    string `json:"slug"`
	Version string `json:"version"`
}

// handleRegistryInstall fetches a plugin from the configured registry,
// verifies its SHA256 against the registry's metadata, and installs
// it via the existing Manager.Install path. Same trust prompts and
// re-approval flow as a manual upload.
func (p *Host) handleRegistryInstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req registryInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: "invalid request body"})
		return
	}
	if req.Slug == "" || req.Version == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: "slug and version required"})
		return
	}

	base := registryBase()
	if base == "" {
		writeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{jsonErrorKey: "plugin registry is not configured on this server"})
		return
	}

	// Fetch the plugin detail from the registry to learn the canonical
	// downloadURL + sha256 for the requested version. We don't trust
	// any URL the admin frontend might have constructed locally; the
	// registry's metadata is the source of truth.
	detail, err := fetchRegistryDetail(base, req.Slug)
	if err != nil {
		writeJSONResponse(w, http.StatusBadGateway, map[string]string{jsonErrorKey: err.Error()})
		return
	}
	match := findVersion(detail, req.Version)
	if match == nil {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{jsonErrorKey: fmt.Sprintf("version %q not found in registry for plugin %q", req.Version, req.Slug)})
		return
	}

	// Pull the bytes from the registry's CDN. LimitReader caps how
	// much we'll read so a hostile registry can't fill memory.
	pkgBytes, err := downloadOcpkg(match.DownloadURL)
	if err != nil {
		writeJSONResponse(w, http.StatusBadGateway, map[string]string{jsonErrorKey: err.Error()})
		return
	}

	// SHA256 over the downloaded bytes must match the registry's
	// stated digest. This is integrity, not authenticity: it tells us
	// the CDN delivered what the registry says is the package, not
	// that the registry itself is trustworthy.
	sum := sha256.Sum256(pkgBytes)
	if hex.EncodeToString(sum[:]) != match.SHA256 {
		writeJSONResponse(w, http.StatusBadGateway, map[string]string{jsonErrorKey: "downloaded package hash does not match registry metadata"})
		return
	}

	entry, err := p.manager.Install(r.Context(), pkgBytes)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: err.Error()})
		return
	}
	log.Infof("plugin %q [%s] v%s installed from registry by admin", entry.DisplayName, entry.Slug, entry.Version)
	writeJSONResponse(w, http.StatusOK, entry)
}

// registryVersion is the subset of the registry's per-version payload
// we need to install a package.
type registryVersion struct {
	Version     string `json:"version"`
	SHA256      string `json:"sha256"`
	DownloadURL string `json:"downloadURL"`
}

type registryDetail struct {
	Slug     string            `json:"slug"`
	Versions []registryVersion `json:"versions"`
}

// fetchRegistryDetail GETs the registry's detail endpoint for a
// single plugin and decodes only the fields we use for install. The
// registry can return extra fields without breaking us. Slug is the
// registry's primary identifier, matching the URL segment used here.
func fetchRegistryDetail(base, slug string) (*registryDetail, error) {
	resp, err := registryHTTPClient.Get(base + "/plugins/" + url.PathEscape(slug))
	if err != nil {
		return nil, fmt.Errorf("registry fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin %q not in registry", slug)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return nil, fmt.Errorf("read registry response: %w", err)
	}
	var detail registryDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("decode registry response: %w", err)
	}
	return &detail, nil
}

func findVersion(detail *registryDetail, version string) *registryVersion {
	for i := range detail.Versions {
		if detail.Versions[i].Version == version {
			return &detail.Versions[i]
		}
	}
	return nil
}

// downloadOcpkg streams the .ocpkg URL the registry pointed at into a
// byte slice, capped at MaxUploadBytes (same limit as direct uploads).
func downloadOcpkg(downloadURL string) ([]byte, error) {
	resp, err := registryHTTPClient.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("download .ocpkg: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download .ocpkg: status %s", resp.Status)
	}
	pkgBytes, err := io.ReadAll(io.LimitReader(resp.Body, plugins.MaxUploadBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read .ocpkg body: %w", err)
	}
	if int64(len(pkgBytes)) > plugins.MaxUploadBytes {
		return nil, fmt.Errorf("downloaded .ocpkg exceeds %d-byte cap", plugins.MaxUploadBytes)
	}
	return pkgBytes, nil
}
