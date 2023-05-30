package storageproviders

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/grafov/m3u8"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/playlist"

	log "github.com/sirupsen/logrus"
)

// rewriteRemotePlaylist will take a local playlist and rewrite it to have absolute URLs to remote locations.
func rewriteRemotePlaylist(localFilePath, remoteServingEndpoint string) error {
	f, err := os.Open(localFilePath) // nolint
	if err != nil {
		log.Fatalln(err)
	}

	p := m3u8.NewMasterPlaylist()
	if err := p.DecodeFrom(bufio.NewReader(f), false); err != nil {
		log.Warnln(err)
	}

	for _, item := range p.Variants {
		item.URI = remoteServingEndpoint + filepath.Join("/hls", item.URI)
	}

	publicPath := filepath.Join(config.HLSStoragePath, filepath.Base(localFilePath))

	newPlaylist := p.String()

	return playlist.WritePlaylist(newPlaylist, publicPath)
}
