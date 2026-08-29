package admin

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncast/owncast/config"
)

// deleteEmojiRequest posts a delete request for the given name and returns
// whether the handler reported success.
func deleteEmojiRequest(t *testing.T, name string) bool {
	t.Helper()

	body, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/admin/emoji/delete", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	DeleteCustomEmoji(rec, req)

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unable to read response %q: %v", rec.Body.String(), err)
	}

	return response.Success
}

// setupEmojiDir points the emoji directory at a temporary working directory
// holding one emoji, the database, and the config file.
func setupEmojiDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.MkdirAll(filepath.Join(dir, config.CustomEmojiPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, config.CustomEmojiPath, "smiley.png"), []byte("emoji"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, config.DataDirectory, "owncast.db"), []byte("database"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("config"), 0o600); err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestDeleteCustomEmojiRemovesEmoji(t *testing.T) {
	dir := setupEmojiDir(t)

	// The admin sends the path with a leading slash.
	if !deleteEmojiRequest(t, "/smiley.png") {
		t.Fatal("expected the emoji to be deleted")
	}

	if _, err := os.Stat(filepath.Join(dir, config.CustomEmojiPath, "smiley.png")); !os.IsNotExist(err) {
		t.Error("emoji was not removed")
	}
}

func TestDeleteCustomEmojiRemovesEmojiInSubdirectory(t *testing.T) {
	dir := setupEmojiDir(t)

	nested := filepath.Join(dir, config.CustomEmojiPath, "pack")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "wave.png"), []byte("emoji"), 0o600); err != nil {
		t.Fatal(err)
	}

	if !deleteEmojiRequest(t, "/pack/wave.png") {
		t.Fatal("expected the nested emoji to be deleted")
	}

	if _, err := os.Stat(filepath.Join(nested, "wave.png")); !os.IsNotExist(err) {
		t.Error("nested emoji was not removed")
	}
}

// Names that leave the emoji directory must be refused. filepath.Join cleans
// ../ away, so checking the joined path let these through (GHSA-2p2g-5h44-ggp2).
func TestDeleteCustomEmojiRefusesPathsOutsideEmojiDirectory(t *testing.T) {
	names := []string{
		"../../config.yaml",
		"../owncast.db",
		"/../../config.yaml",
		"..",
		".",
		"",
	}

	for i, name := range names {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Helper()
			dir := setupEmojiDir(t)

			if deleteEmojiRequest(t, name) {
				t.Errorf("name %q was accepted", name)
			}

			for _, survivor := range []string{"config.yaml", filepath.Join(config.DataDirectory, "owncast.db")} {
				if _, err := os.Stat(filepath.Join(dir, survivor)); err != nil {
					t.Errorf("%s outside of the emoji directory was removed: %v", survivor, err)
				}
			}
			if _, err := os.Stat(filepath.Join(dir, config.CustomEmojiPath)); err != nil {
				t.Errorf("emoji directory was removed: %v", err)
			}
		})
	}
}

// A symlink inside the emoji directory must not become a way out of it.
func TestDeleteCustomEmojiRefusesSymlinkedPathOutsideEmojiDirectory(t *testing.T) {
	dir := setupEmojiDir(t)

	outside := filepath.Join(dir, "outside")
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(outside, "secret.txt")
	if err := os.WriteFile(secret, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, config.CustomEmojiPath, "link")); err != nil {
		t.Fatal(err)
	}

	if deleteEmojiRequest(t, "/link/secret.txt") {
		t.Error("a path through a symlink was accepted")
	}

	if _, err := os.Stat(secret); err != nil {
		t.Errorf("file outside of the emoji directory was removed: %v", err)
	}
}
