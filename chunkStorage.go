package main

type ChunkStorage interface {
	Setup(config Config)
	Save(filePath string) string
	GenerateRemotePlaylist(playlist string, segments map[string]string) string
}
