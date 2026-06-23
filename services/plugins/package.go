package plugins

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// Decompressed-size caps for the entries the host reads out of an .ocpkg. The
// archive itself is capped (MaxUploadBytes), but a small compressed entry can
// inflate to gigabytes (a zip bomb) and OOM the host, so every read of a
// decompressed entry is bounded. The manifest cap is the tightest because the
// manifest is re-read on every plugins-directory scan.
const (
	maxManifestEntryBytes     = 4 << 20   // 4 MiB
	maxCodeEntryBytes         = 128 << 20 // 128 MiB (self-contained wasm can be large)
	maxIconEntryBytes         = 8 << 20   // 8 MiB
	maxInstructionsEntryBytes = 2 << 20   // 2 MiB
)

// The .ocpkg ("Owncast plugin package") format is a zip archive with four
// well-known entries:
//
//	plugin.manifest.json    required, the sidecar manifest
//	plugin.wasm             required, the compiled plugin module
//	public/...              optional, files served at /plugins/<name>/<path>
//	assets/...              optional, files the host reads internally for
//	                        manifest fields that inline content
//	                        (styles, scripts, extraPageContent); never
//	                        reachable through the plugin's URL space.
//
// File names inside the archive are canonical regardless of the plugin's name
// so the host doesn't have to read the manifest to discover the wasm path.
// The plugin's name still comes from manifest.name, not the .ocpkg filename.
const (
	packageSuffix       = ".ocpkg"
	pkgManifestFilename = "plugin.manifest.json"
	// Code entries. The plugin's runtime is inferred from which one is present
	// — the author doesn't declare it in the manifest. plugin.js / plugin.py
	// are author source run on the shared embedded engine; plugin.wasm is a
	// self-contained module the host loads directly.
	pkgWasmFilename = "plugin.wasm"
	pkgJSFilename   = "plugin.js"
	pkgPyFilename   = "plugin.py"
	pkgPublicPrefix = "public/"
	pkgAssetsPrefix = "assets/"
)

// codeEntries maps each canonical code-entry name to the runtime it implies,
// in detection-precedence order.
var codeEntries = []struct{ name, runtime string }{
	{pkgJSFilename, RuntimeJavaScript},
	{pkgPyFilename, RuntimePython},
	{pkgWasmFilename, RuntimeWasm},
}

// detectPackageCode finds a package's code entry and the runtime it implies.
func detectPackageCode(zr *zip.Reader) (name, runtime string, ok bool) {
	for _, e := range codeEntries {
		if zipHasFile(zr, e.name) {
			return e.name, e.runtime, true
		}
	}
	return "", "", false
}

func zipHasFile(zr *zip.Reader, name string) bool {
	for _, f := range zr.File {
		if f.Name == name {
			return true
		}
	}
	return false
}

// LoadPackage loads a plugin from a .ocpkg file. The archive is opened with
// a file-backed reader so only the central directory and the manifest/wasm
// entries actually enter memory; per-asset reads happen on demand when the
// HTTP server fetches them. The *zip.ReadCloser is retained on the returned
// Loaded so the underlying file stays open for AssetsFS lookups, and is
// closed when Loaded.Close runs.
func LoadPackage(ctx context.Context, env *HostEnv, path string) (*Loaded, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open %s as zip: %w", path, err)
	}
	// Until we hand ownership to the returned Loaded, any error path must
	// close zr or we leak the file handle.
	closeOnFail := zr
	defer func() {
		if closeOnFail != nil {
			_ = closeOnFail.Close()
		}
	}()

	manifestBytes, err := readZipFile(&zr.Reader, pkgManifestFilename, maxManifestEntryBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	// The code entry's name implies the runtime (plugin.js / plugin.py / plugin.wasm),
	// so the author never declares it in the manifest.
	codeName, runtimeType, ok := detectPackageCode(&zr.Reader)
	if !ok {
		return nil, fmt.Errorf("%s: missing plugin code (expected one of plugin.js, plugin.py, plugin.wasm)", filepath.Base(path))
	}
	codeBytes, err := readZipFile(&zr.Reader, codeName, maxCodeEntryBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}

	displayName := strings.TrimSuffix(filepath.Base(path), packageSuffix)

	// Extract assetsFS from the zip before calling loadFromBytes so the
	// owncast_asset_read host function has access to it at instantiation time.
	var assetsFS fs.FS
	if hasZipDir(&zr.Reader, pkgAssetsPrefix) {
		if sub, err := fs.Sub(&zr.Reader, strings.TrimSuffix(pkgAssetsPrefix, "/")); err == nil {
			assetsFS = sub
		}
	}

	loaded, err := loadFromBytes(ctx, env, manifestBytes, codeBytes, runtimeType, displayName, assetsFS)
	if err != nil {
		return nil, err
	}
	loaded.WasmPath = path
	loaded.pkgCloser = zr

	// Mount public/ as the plugin's web-served root. AssetsFS is already set
	// by loadFromBytes. fs.Sub returns an FS that's empty (rather than failing)
	// when the prefix doesn't exist, so we check first to keep the
	// nil-means-empty invariant the Server (PublicFS) relies on.
	if hasZipDir(&zr.Reader, pkgPublicPrefix) {
		if sub, err := fs.Sub(&zr.Reader, strings.TrimSuffix(pkgPublicPrefix, "/")); err == nil {
			loaded.PublicFS = sub
		}
	}

	closeOnFail = nil // ownership transferred to Loaded.pkgCloser
	return loaded, nil
}

// readManifestFromPackage reads just the plugin.manifest.json entry from a
// .ocpkg file without instantiating the wasm. Used by Manager.scan() — runs
// every couple seconds, so the file-backed reader matters: a multi-gigabyte
// .ocpkg sitting in plugins/ costs only a central-directory read per scan,
// not a full slurp into memory.
func readManifestFromPackage(path string) (*Manifest, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open %s as zip: %w", path, err)
	}
	defer zr.Close()
	manifestBytes, err := readZipFile(&zr.Reader, pkgManifestFilename, maxManifestEntryBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	return ParseManifest(manifestBytes)
}

func readZipFile(zr *zip.Reader, name string, maxBytes int64) ([]byte, error) {
	for _, f := range zr.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("open %s: %w", name, err)
			}
			defer rc.Close()
			// Bound the DECOMPRESSED size: read one byte past the cap so we can
			// tell "exactly at the cap" from "over it" without trusting the
			// zip header's claimed size.
			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(io.LimitReader(rc, maxBytes+1)); err != nil {
				return nil, fmt.Errorf("read %s: %w", name, err)
			}
			if int64(buf.Len()) > maxBytes {
				return nil, fmt.Errorf("%s exceeds the %d-byte limit", name, maxBytes)
			}
			return buf.Bytes(), nil
		}
	}
	return nil, fmt.Errorf("missing required entry %q", name)
}

func hasZipDir(zr *zip.Reader, prefix string) bool {
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, prefix) {
			return true
		}
	}
	return false
}
