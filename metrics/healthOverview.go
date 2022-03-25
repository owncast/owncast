package metrics

import (
	"fmt"
	"sort"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

var errorMessages = map[string]string{
	"LOWSPEED":        "%d of %d clients (%d%%) are consuming video slower than, or too close to your bitrate of %d kbps.",
	"PLAYBACK_ERRORS": "%d of %d clients (%d%%) are experiencing different, unspecified, playback issues.",
}

// GetStreamHealthOverview will return the stream health overview.
func GetStreamHealthOverview() *models.StreamHealthOverview {
	return metrics.streamHealthOverview
}

func generateStreamHealthOverview() {
	overview := models.StreamHealthOverview{
		Healthy:           true,
		HealthyPercentage: 100,
	}

	defer func() {
		metrics.streamHealthOverview = &overview
	}()

	type singleVariant struct {
		isVideoPassthrough bool
		bitrate            int
	}

	outputVariants := data.GetStreamOutputVariants()

	streamSortVariants := make([]singleVariant, len(outputVariants))
	for i, variant := range outputVariants {
		variantSort := singleVariant{
			bitrate:            variant.VideoBitrate,
			isVideoPassthrough: variant.IsVideoPassthrough,
		}
		streamSortVariants[i] = variantSort
	}

	sort.Slice(streamSortVariants, func(i, j int) bool {
		if streamSortVariants[i].isVideoPassthrough && !streamSortVariants[j].isVideoPassthrough {
			return true
		}

		if !streamSortVariants[i].isVideoPassthrough && streamSortVariants[j].isVideoPassthrough {
			return false
		}

		return streamSortVariants[i].bitrate > streamSortVariants[j].bitrate
	})

	lowestSupportedBitrate := float64(streamSortVariants[0].bitrate)
	totalNumberOfClients := len(windowedBandwidths)

	if totalNumberOfClients == 0 {
		return
	}

	// Determine healthy status based on bandwidth speeds of clients.
	unhealthyClientCount := 0
	for _, speed := range windowedBandwidths {
		if int(speed) < int(lowestSupportedBitrate*1.1) {
			unhealthyClientCount++
		}
	}
	if unhealthyClientCount > 0 {
		overview.Message = fmt.Sprintf(errorMessages["LOWSPEED"], unhealthyClientCount, totalNumberOfClients, int((float64(unhealthyClientCount)/float64(totalNumberOfClients))*100), int(lowestSupportedBitrate))
	}

	// If bandwidth is ok, determine healthy status based on error
	// counts of clients.
	if unhealthyClientCount == 0 {
		for _, errors := range windowedErrorCounts {
			unhealthyClientCount += int(errors)
		}
		if unhealthyClientCount > 0 {
			overview.Message = fmt.Sprintf(errorMessages["PLAYBACK_ERRORS"], unhealthyClientCount, totalNumberOfClients, int((float64(unhealthyClientCount)/float64(totalNumberOfClients))*100))
		}
	}

	if unhealthyClientCount == 0 {
		return
	}

	percentUnhealthy := 100 - ((float64(unhealthyClientCount) / float64(totalNumberOfClients)) * 100)
	overview.HealthyPercentage = int(percentUnhealthy)
	overview.Healthy = overview.HealthyPercentage > 95
}
