package replays

type StorageProvider interface {
	Setup() error
	Save(localFilePath, destinationPath string, retryCount int) (string, error)
}
