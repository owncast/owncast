package models

// StorageProvider is how a chunk storage provider should be implemented.
type StorageProvider interface {
	Setup() error
	Save(localFilePath, destinationPath string, retryCount int) (string, error)
	SetStreamId(streamID string)
	SegmentWritten(localFilePath string) (string, int, error)
	VariantPlaylistWritten(localFilePath string)
	MasterPlaylistWritten(localFilePath string)
	GetRemoteDestinationPathFromLocalFilePath(localFilePath string) string

	Cleanup() error
}
