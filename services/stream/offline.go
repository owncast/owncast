package stream

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
)

func (s *Service) appendOfflineToVariantPlaylist(index int, playlistFilePath string) {
	existingPlaylistContents, err := os.ReadFile(playlistFilePath) // nolint: gosec
	if err != nil {
		log.Debugln("unable to read existing playlist file", err)
		return
	}

	tmpFileName := fmt.Sprintf("tmp-stream-%d.m3u8", index)
	atomicWriteTmpPlaylistFile, err := os.CreateTemp(s.cfg.TempDir, tmpFileName)
	if err != nil {
		log.Errorln("error creating tmp playlist file to write to", playlistFilePath, err)
		return
	}

	// Write the existing playlist contents
	if _, err := atomicWriteTmpPlaylistFile.Write(existingPlaylistContents); err != nil {
		log.Debugln("error writing existing playlist contents to tmp playlist file", err)
		return
	}

	// Manually append the offline fMP4 clip to the end of the media playlist.
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXT-X-DISCONTINUITY\n")
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXT-X-MAP:URI=\"offline-init.mp4\"\n")
	// If "offline" content gets changed then change the duration below
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXTINF:8.000000,\n")
	_, _ = atomicWriteTmpPlaylistFile.WriteString("offline-v2.m4s\n")
	_, _ = atomicWriteTmpPlaylistFile.WriteString("#EXT-X-ENDLIST\n")

	if err := atomicWriteTmpPlaylistFile.Close(); err != nil {
		log.Errorln(err)
	}

	if err := utils.Move(atomicWriteTmpPlaylistFile.Name(), playlistFilePath); err != nil {
		log.Errorln("error moving temp playlist to overwrite existing one", err)
	}
}

func (s *Service) makeVariantIndexOffline(index int, offlineInitSrcPath string, offlineSegmentSrcPath string) {
	playlistFilePath := fmt.Sprintf(filepath.Join(config.HLSStoragePath, "%d/stream.m3u8"), index)

	// Copy the fMP4 init segment for offline playback.
	initDestPath := fmt.Sprintf(filepath.Join(config.HLSStoragePath, "%d/offline-init.mp4"), index)
	if err := utils.Copy(offlineInitSrcPath, initDestPath); err != nil {
		log.Warnln(err)
	}
	if _, err := s.storage.Save(initDestPath, 0); err != nil {
		log.Warnln(err)
	}

	// Copy the fMP4 media segment for offline playback.
	segmentDestPath := fmt.Sprintf(filepath.Join(config.HLSStoragePath, "%d/offline-v2.m4s"), index)
	if err := utils.Copy(offlineSegmentSrcPath, segmentDestPath); err != nil {
		log.Warnln(err)
	}
	if _, err := s.storage.Save(segmentDestPath, 0); err != nil {
		log.Warnln(err)
	}

	if utils.DoesFileExists(playlistFilePath) {
		s.appendOfflineToVariantPlaylist(index, playlistFilePath)
	} else {
		createEmptyOfflinePlaylist(playlistFilePath)
	}
	if _, err := s.storage.Save(playlistFilePath, 0); err != nil {
		log.Warnln(err)
	}
}

func createEmptyOfflinePlaylist(playlistFilePath string) {
	f, err := os.Create(playlistFilePath) //nolint:gosec
	if err != nil {
		log.Errorln(err)
		return
	}

	// Write a minimal fMP4-compatible HLS media playlist for the offline state.
	content := "#EXTM3U\n" +
		"#EXT-X-VERSION:6\n" +
		"#EXT-X-TARGETDURATION:8\n" +
		"#EXT-X-MEDIA-SEQUENCE:0\n" +
		"#EXT-X-MAP:URI=\"offline-init.mp4\"\n" +
		// If "offline" content gets changed then change the duration below
		"#EXTINF:8.000000,\n" +
		"offline-v2.m4s\n" +
		"#EXT-X-ENDLIST\n"

	if _, err = f.WriteString(content); err != nil {
		log.Errorln("error writing empty offline playlist:", err)
	}

	if err = f.Close(); err != nil {
		log.Errorln("error closing offline playlist file:", err)
	}
}

func saveOfflineFMP4ToDisk(tempDir string) (initPath string, segmentPath string, err error) {
	initData := static.GetOfflineInitSegment()
	initTmp, err := os.CreateTemp(tempDir, "offline-init-*.mp4")
	if err != nil {
		return "", "", fmt.Errorf("unable to create temp file for offline init segment: %s", err)
	}

	if _, err = initTmp.Write(initData); err != nil {
		return "", "", fmt.Errorf("unable to write offline init segment to disk: %s", err)
	}

	if err = initTmp.Close(); err != nil {
		return "", "", fmt.Errorf("unable to close offline init segment temp file: %s", err)
	}
	initPath, _ = filepath.Abs(initTmp.Name())

	segData := static.GetOfflineMediaSegment()
	segTmp, err := os.CreateTemp(tempDir, "offline-v2-*.m4s")
	if err != nil {
		return "", "", fmt.Errorf("unable to create temp file for offline media segment: %s", err)
	}

	if _, err = segTmp.Write(segData); err != nil {
		return "", "", fmt.Errorf("unable to write offline media segment to disk: %s", err)
	}

	if err = segTmp.Close(); err != nil {
		return "", "", fmt.Errorf("unable to close offline media segment temp file: %s", err)
	}
	segmentPath, _ = filepath.Abs(segTmp.Name())

	return initPath, segmentPath, nil
}

func saveOfflineClipToDisk(tempDir, offlineFilename string) (string, error) {
	offlineFileData := static.GetOfflineSegment()
	offlineTmpFile, err := os.CreateTemp(tempDir, offlineFilename)
	if err != nil {
		log.Errorln("unable to create temp file for offline video segment", err)
	}

	if _, err = offlineTmpFile.Write(offlineFileData); err != nil {
		return "", fmt.Errorf("unable to write offline segment to disk: %s", err)
	}

	offlineFilePath := offlineTmpFile.Name()
	return filepath.Abs(offlineFilePath)
}
