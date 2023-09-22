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

// rewritePlaylistLocations will take a local playlist and rewrite it to have absolute URLs to a specified location.
func rewritePlaylistLocations(localFilePath, remoteServingEndpoint, pathPrefix string) error {
	f, err := os.Open(localFilePath) // nolint
	if err != nil {
		log.Fatalln(err)
	}

	p := m3u8.NewMasterPlaylist()
	if err := p.DecodeFrom(bufio.NewReader(f), false); err != nil {
		log.Warnln(err)
	}

	for _, item := range p.Variants {
		// Determine the final path to this playlist.
		var finalPath string
		if pathPrefix != "" {
			finalPath = filepath.Join(pathPrefix, "/hls")
		} else {
			finalPath = "/hls"
		}
		item.URI = remoteServingEndpoint + filepath.Join(finalPath, item.URI)
	}

	publicPath := filepath.Join(config.HLSStoragePath, filepath.Base(localFilePath))

	newPlaylist := p.String()

	return playlist.WritePlaylist(newPlaylist, publicPath)
}
