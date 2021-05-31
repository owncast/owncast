package rtmp

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nareix/joy5/format/flv/flvio"
	"github.com/owncast/owncast/models"
)

const unknownString = "Unknown"

func getInboundDetailsFromMetadata(metadata []interface{}) (models.RTMPStreamMetadata, error) {
	metadataComponentsString := fmt.Sprintf("%+v", metadata)
	if !strings.Contains(metadataComponentsString, "onMetaData") {
		return models.RTMPStreamMetadata{}, errors.New("Not a onMetaData message")
	}
	re := regexp.MustCompile(`\{(.*?)\}`)
	submatchall := re.FindAllString(metadataComponentsString, 1)

	if len(submatchall) == 0 {
		return models.RTMPStreamMetadata{}, errors.New("unable to parse inbound metadata")
	}

	metadataJSONString := submatchall[0]
	var details models.RTMPStreamMetadata
	err := json.Unmarshal([]byte(metadataJSONString), &details)
	return details, err
}

func getAudioCodec(codec interface{}) string {
	if codec == nil {
		return "No audio"
	}

	var codecID float64
	if assertedCodecID, ok := codec.(float64); ok {
		codecID = assertedCodecID
	} else {
		return codec.(string)
	}

	switch codecID {
	case flvio.SOUND_MP3:
		return "MP3"
	case flvio.SOUND_AAC:
		return "AAC"
	case flvio.SOUND_SPEEX:
		return "Speex"
	}

	return unknownString
}

func getVideoCodec(codec interface{}) string {
	if codec == nil {
		return unknownString
	}

	var codecID float64
	if assertedCodecID, ok := codec.(float64); ok {
		codecID = assertedCodecID
	} else {
		return codec.(string)
	}

	switch codecID {
	case flvio.VIDEO_H264:
		return "H.264"
	case flvio.VIDEO_H265:
		return "H.265"
	}

	return unknownString
}

func secretMatch(configStreamKey string, url string) bool {
  sep := "/live/"
  givenParts := strings.Split(url, sep)

  if len(givenParts) < 2 { // url part and secret
    return false
  }

  // If the given url has /live/ in the secret, this /live/ was removed with the
  // `strings.Split/2` above. So, we add it back now.
  streamingKey := strings.Join(givenParts[1:], sep)

  return streamingKey == configStreamKey
}
