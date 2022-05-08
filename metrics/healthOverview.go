package metrics

import (
	"fmt"
	"sort"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

const (
	healthyPercentageMinValue = 75
	maxCPUUsage               = 90
	minClientCountForDetails  = 3
)

// GetStreamHealthOverview will return the stream health overview.
func GetStreamHealthOverview() *models.StreamHealthOverview {
	return metrics.streamHealthOverview
}

func generateStreamHealthOverview() {
	// Determine what percentage of total players are represented in our overview.
	totalPlayerCount := len(core.GetActiveViewers())
	if totalPlayerCount == 0 {
		metrics.streamHealthOverview = nil
		return
	}

	pct := getClientErrorHeathyPercentage()
	if pct < 1 {
		metrics.streamHealthOverview = nil
		return
	}

	overview := &models.StreamHealthOverview{
		Healthy:           pct > healthyPercentageMinValue,
		HealthyPercentage: pct,
		Message:           getStreamHealthOverviewMessage(),
	}

	if totalPlayerCount > 0 && len(windowedBandwidths) > 0 {
		representation := utils.IntPercentage(len(windowedBandwidths), totalPlayerCount)
		overview.Representation = representation
	}

	metrics.streamHealthOverview = overview
}

func getStreamHealthOverviewMessage() string {
	if message := wastefulBitrateOverviewMessage(); message != "" {
		return message
	} else if message := cpuUsageHealthOverviewMessage(); message != "" {
		return message
	} else if message := networkSpeedHealthOverviewMessage(); message != "" {
		return message
	} else if message := errorCountHealthOverviewMessage(); message != "" {
		return message
	}

	return ""
}

func networkSpeedHealthOverviewMessage() string {
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

	lowestSupportedBitrate := float64(streamSortVariants[len(streamSortVariants)-1].bitrate)
	totalNumberOfClients := len(windowedBandwidths)

	if totalNumberOfClients == 0 {
		return ""
	}

	// Determine healthy status based on bandwidth speeds of clients.
	unhealthyClientCount := 0

	for _, speed := range windowedBandwidths {
		if int(speed) < int(lowestSupportedBitrate*1.1) {
			unhealthyClientCount++
		}
	}

	if unhealthyClientCount == 0 {
		return ""
	}

	return fmt.Sprintf("%d of %d viewers (%d%%) are consuming video slower than, or too close to your bitrate of %d kbps.", unhealthyClientCount, totalNumberOfClients, int((float64(unhealthyClientCount)/float64(totalNumberOfClients))*100), int(lowestSupportedBitrate))
}

// wastefulBitrateOverviewMessage attempts to determine if a streamer is sending to
// Owncast at a bitrate higher than they're streaming to their viewers leading
// to wasted CPU by having to compress it.
func wastefulBitrateOverviewMessage() string {
	if len(metrics.CPUUtilizations) < 2 {
		return ""
	}

	// Only return an alert if the CPU usage is around the max cpu threshold.
	recentCPUUses := metrics.CPUUtilizations[len(metrics.CPUUtilizations)-2:]
	values := make([]float64, len(recentCPUUses))
	for i, val := range recentCPUUses {
		values[i] = val.Value
	}
	recentCPUUse := utils.Avg(values)

	if recentCPUUse < maxCPUUsage-10 {
		return ""
	}

	currentBroadcast := core.GetCurrentBroadcast()
	if currentBroadcast == nil {
		return ""
	}

	currentBroadcaster := core.GetBroadcaster()
	if currentBroadcast == nil {
		return ""
	}

	if currentBroadcaster.StreamDetails.AudioBitrate == 0 {
		return ""
	}

	inboundBitrate := currentBroadcaster.StreamDetails.VideoBitrate
	maxBitrate := 0
	if inboundBitrate > maxBitrate {
		return fmt.Sprintf("You're broadcasting to Owncast at %dkbps but only sending to your viewers at %dkbps, requiring unnecessary work to be performed. You may want to decrease what you're sending to Owncast or increase what you send to your viewers to match.", inboundBitrate, maxBitrate)
	}

	return ""
}

func cpuUsageHealthOverviewMessage() string {
	if len(metrics.CPUUtilizations) < 2 {
		return ""
	}

	recentCPUUses := metrics.CPUUtilizations[len(metrics.CPUUtilizations)-2:]
	values := make([]float64, len(recentCPUUses))
	for i, val := range recentCPUUses {
		values[i] = val.Value
	}
	recentCPUUse := utils.Avg(values)
	if recentCPUUse < maxCPUUsage {
		return ""
	}

	return fmt.Sprintf("The CPU usage on your server is over %d%%. This may cause video to be provided slower than necessary, causing buffering for your viewers. Consider increasing the resources available or reducing the number of output variants you made available.", maxCPUUsage)
}

func errorCountHealthOverviewMessage() string {
	totalNumberOfClients := len(windowedBandwidths)
	if totalNumberOfClients == 0 {
		return ""
	}

	clientsWithErrors := getClientsWithErrorsCount()

	if clientsWithErrors == 0 {
		return ""
	}

	// Only return these detailed values and messages if we feel we have enough
	// clients to be able to make a reasonable assessment. This is an arbitrary
	// number but 1 out of 1 isn't helpful.

	if totalNumberOfClients >= minClientCountForDetails {
		healthyPercentage := utils.IntPercentage(clientsWithErrors, totalNumberOfClients)

		isUsingPassthrough := false
		outputVariants := data.GetStreamOutputVariants()
		for _, variant := range outputVariants {
			if variant.IsVideoPassthrough {
				isUsingPassthrough = true
			}
		}

		if isUsingPassthrough {
			return fmt.Sprintf("%d of %d viewers (%d%%) are experiencing errors. You're currently using a video passthrough output, often known for causing playback issues for people. It is suggested you turn it off.", clientsWithErrors, totalNumberOfClients, healthyPercentage)
		}

		currentBroadcast := core.GetCurrentBroadcast()
		if currentBroadcast != nil && currentBroadcast.LatencyLevel.SecondsPerSegment < 3 {
			return fmt.Sprintf("%d of %d viewers (%d%%) may be experiencing some issues. You may want to increase your latency buffer level in your video configuration to see if it helps.", clientsWithErrors, totalNumberOfClients, healthyPercentage)
		}

		return fmt.Sprintf("%d of %d viewers (%d%%) may be experiencing some issues.", clientsWithErrors, totalNumberOfClients, healthyPercentage)
	}

	return ""
}

func getClientsWithErrorsCount() int {
	clientsWithErrors := 0
	for _, errors := range windowedErrorCounts {
		if errors > 0 {
			clientsWithErrors++
		}
	}
	return clientsWithErrors
}

func getClientErrorHeathyPercentage() int {
	totalNumberOfClients := len(windowedErrorCounts)
	if totalNumberOfClients == 0 {
		return -1
	}

	clientsWithErrors := getClientsWithErrorsCount()

	if clientsWithErrors == 0 {
		return 100
	}

	pct := 100 - utils.IntPercentage(clientsWithErrors, totalNumberOfClients)

	return pct
}
