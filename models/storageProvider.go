package models

//ChunkStorageProvider is how a chunk storage provider should be implemented
type ChunkStorageProvider interface {
	Setup() error
	Save(filePath string, retryCount int) (string, error)
	GenerateRemotePlaylist(playlist string, variant Variant) string
}
