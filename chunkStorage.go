package main

type ChunkStorage interface {
	Setup(config Config)
	Save(filePath string, retryCount int) string
	GenerateRemotePlaylist(playlist string, variant Variant) string
}
