package models

// StorageProvider is how a chunk storage provider should be implemented.
type StorageProvider interface {
	Setup() error
	Save(filePath string, retryCount int) (string, error)

	// SegmentWritten is called when a single HLS segment is written. It
	// returns the public path the segment is served from (a relative path
	// for local storage, an absolute URL for remote storage) so the replay
	// recorder can reference it later.
	SegmentWritten(localFilePath string) (string, error)
	VariantPlaylistWritten(localFilePath string)
	MasterPlaylistWritten(localFilePath string)

	Cleanup() error
}
