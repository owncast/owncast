package rtmp

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/owncast/owncast/models"
	"github.com/nareix/joy5/format/flv/flvio"
)

func getInboundDetailsFromMetadata(metadata []interface{}) (models.RTMPStreamMetadata, error) {
	metadataComponentsString := fmt.Sprintf("%+v", metadata)
	re := regexp.MustCompile(`\{(.*?)\}`)
	submatchall := re.FindAllString(metadataComponentsString, 1)

	if len(submatchall) == 0 {
		return models.RTMPStreamMetadata{}, errors.New("unable to parse inbound metadata")
	}

	metadataJSONString := submatchall[0]
	var details models.RTMPStreamMetadata
	json.Unmarshal([]byte(metadataJSONString), &details)
	return details, nil
}

func getAudioCodec(codec interface{}) string {
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

	return "Unknown"
}

func getVideoCodec(codec interface{}) string {
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

	return "Unknown"
}
