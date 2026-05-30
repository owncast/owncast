package storage

import "os"

// writePlaylist writes an HLS playlist to disk. Internal helper for the
// storage providers (local + S3 both rely on rewriting playlist files
// to point at their final serving location).
func writePlaylist(data string, filePath string) error {
	// nolint:gosec
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return nil
}
