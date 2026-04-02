package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

func appendOfflineToVariantPlaylist(index int, playlistFilePath string, initFilename string, segmentFilename string) {
	existingPlaylistContents, err := os.ReadFile(playlistFilePath) // nolint: gosec
	if err != nil {
		log.Debugln("unable to read existing playlist file", err)
		return
	}

	tmpFileName := fmt.Sprintf("tmp-stream-%d.m3u8", index)
	atomicWriteTmpPlaylistFile, err := os.CreateTemp(config.TempDir, tmpFileName)
	if err != nil {
		log.Errorln("error creating tmp playlist file to write to", playlistFilePath, err)
		return
	}

	// Write the existing playlist contents
	if _, err := atomicWriteTmpPlaylistFile.Write(existingPlaylistContents); err != nil {
		log.Debugln("error writing existing playlist contents to tmp playlist file", err)
		return
	}

	// Manually append the offline clip to the end of the media playlist.
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXT-X-DISCONTINUITY\n")
	_, _ = fmt.Fprintf(atomicWriteTmpPlaylistFile, "#EXT-X-MAP:URI=\"%s\"\n", initFilename)
	// If "offline" content gets changed then change the duration below
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXTINF:8.000000,\n")
	_, _ = fmt.Fprintf(atomicWriteTmpPlaylistFile, "%s\n", segmentFilename)
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXT-X-ENDLIST\n")

	if err := atomicWriteTmpPlaylistFile.Close(); err != nil {
		log.Errorln(err)
	}

	if err := utils.Move(atomicWriteTmpPlaylistFile.Name(), playlistFilePath); err != nil {
		log.Errorln("error moving temp playlist to overwrite existing one", err)
	}
}

func makeVariantIndexOffline(index int, offlineInitPath string, offlineSegmentPath string) {
	variantDir := filepath.Join(config.HLSStoragePath, fmt.Sprintf("%d", index))
	playlistFilePath := filepath.Join(variantDir, "stream.m3u8")

	initDest := filepath.Join(variantDir, "offline-init.mp4")
	segmentDest := filepath.Join(variantDir, "offline-v2.m4s")

	if err := utils.Copy(offlineInitPath, initDest); err != nil {
		log.Warnln(err)
		return
	}

	if _, err := _storage.Save(initDest, 0); err != nil {
		log.Warnln(err)
	}

	if err := utils.Copy(offlineSegmentPath, segmentDest); err != nil {
		log.Warnln(err)
		return
	}

	if _, err := _storage.Save(segmentDest, 0); err != nil {
		log.Warnln(err)
	}

	if utils.DoesFileExists(playlistFilePath) {
		appendOfflineToVariantPlaylist(index, playlistFilePath, "offline-init.mp4", "offline-v2.m4s")
	} else {
		createEmptyOfflinePlaylist(playlistFilePath, "offline-init.mp4", "offline-v2.m4s")
	}
	if _, err := _storage.Save(playlistFilePath, 0); err != nil {
		log.Warnln(err)
	}
}

func createEmptyOfflinePlaylist(playlistFilePath string, initFilename string, segmentFilename string) {
	f, err := os.Create(playlistFilePath) //nolint:gosec
	if err != nil {
		log.Errorln(err)
		return
	}
	defer f.Close()

	_, _ = f.WriteString("#EXTM3U\n")
	_, _ = f.WriteString("#EXT-X-VERSION:7\n")
	_, _ = f.WriteString("#EXT-X-TARGETDURATION:8\n")
	_, _ = fmt.Fprintf(f, "#EXT-X-MAP:URI=\"%s\"\n", initFilename)
	// If "offline" content gets changed then change the duration below
	_, _ = f.WriteString("#EXTINF:8.000000,\n")
	_, _ = fmt.Fprintf(f, "%s\n", segmentFilename)
	_, _ = f.WriteString("#EXT-X-ENDLIST\n")
}

func saveOfflineFMP4ToDisk() (initPath string, segmentPath string, err error) {
	initData := static.GetOfflineInitSegment()
	initTmp, err := os.CreateTemp(config.TempDir, "offline-init-*.mp4")
	if err != nil {
		return "", "", fmt.Errorf("unable to create temp file for offline init segment: %s", err)
	}

	if _, err = initTmp.Write(initData); err != nil {
		return "", "", fmt.Errorf("unable to write offline init segment to disk: %s", err)
	}

	initTmp.Close()
	initPath, _ = filepath.Abs(initTmp.Name())

	segData := static.GetOfflineMediaSegment()
	segTmp, err := os.CreateTemp(config.TempDir, "offline-v2-*.m4s")
	if err != nil {
		return "", "", fmt.Errorf("unable to create temp file for offline media segment: %s", err)
	}

	if _, err = segTmp.Write(segData); err != nil {
		return "", "", fmt.Errorf("unable to write offline media segment to disk: %s", err)
	}

	segTmp.Close()
	segmentPath, _ = filepath.Abs(segTmp.Name())

	return initPath, segmentPath, nil
}
