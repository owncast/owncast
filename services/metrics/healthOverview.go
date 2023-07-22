package metrics

import (
	"fmt"
	"sort"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/status"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/utils"
)

const (
	healthyPercentageMinValue = 75
	maxCPUUsage               = 90
	minClientCountForDetails  = 3
)

// GetStreamHealthOverview will return the stream health overview.
func (m *Metrics) GetStreamHealthOverview() *models.StreamHealthOverview {
	return m.metrics.streamHealthOverview
}

func (m *Metrics) generateStreamHealthOverview() {
	// Determine what percentage of total players are represented in our overview.
	totalPlayerCount := len(core.GetActiveViewers())
	if totalPlayerCount == 0 {
		m.metrics.streamHealthOverview = nil
		return
	}

	pct := m.getClientErrorHeathyPercentage()
	if pct < 1 {
		m.metrics.streamHealthOverview = nil
		return
	}

	overview := &models.StreamHealthOverview{
		Healthy:           pct > healthyPercentageMinValue,
		HealthyPercentage: pct,
		Message:           m.getStreamHealthOverviewMessage(),
	}

	if totalPlayerCount > 0 && len(m.windowedBandwidths) > 0 {
		representation := utils.IntPercentage(len(m.windowedBandwidths), totalPlayerCount)
		overview.Representation = representation
	}

	m.metrics.streamHealthOverview = overview
}

func (m *Metrics) getStreamHealthOverviewMessage() string {
	if message := m.wastefulBitrateOverviewMessage(); message != "" {
		return message
	} else if message := m.cpuUsageHealthOverviewMessage(); message != "" {
		return message
	} else if message := m.networkSpeedHealthOverviewMessage(); message != "" {
		return message
	} else if message := m.errorCountHealthOverviewMessage(); message != "" {
		return message
	}

	return ""
}

func (m *Metrics) networkSpeedHealthOverviewMessage() string {
	type singleVariant struct {
		isVideoPassthrough bool
		bitrate            int
	}

	configRepository := configrepository.Get()

	outputVariants := configRepository.GetStreamOutputVariants()

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
	totalNumberOfClients := len(m.windowedBandwidths)

	if totalNumberOfClients == 0 {
		return ""
	}

	// Determine healthy status based on bandwidth speeds of clients.
	unhealthyClientCount := 0

	for _, speed := range m.windowedBandwidths {
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
func (m *Metrics) wastefulBitrateOverviewMessage() string {
	if len(m.metrics.CPUUtilizations) < 2 {
		return ""
	}

	// Only return an alert if the CPU usage is around the max cpu threshold.
	recentCPUUses := m.metrics.CPUUtilizations[len(m.metrics.CPUUtilizations)-2:]
	values := make([]float64, len(recentCPUUses))
	for i, val := range recentCPUUses {
		values[i] = val.Value
	}
	recentCPUUse := utils.Avg(values)

	if recentCPUUse < maxCPUUsage-10 {
		return ""
	}

	stat := status.Get()

	currentBroadcast := stat.GetCurrentBroadcast()
	if currentBroadcast == nil {
		return ""
	}

	currentBroadcaster := stat.GetBroadcaster()
	if currentBroadcast == nil {
		return ""
	}

	if currentBroadcaster.StreamDetails.VideoBitrate == 0 {
		return ""
	}

	// Not all streams report their inbound bitrate.
	inboundBitrate := currentBroadcaster.StreamDetails.VideoBitrate
	if inboundBitrate == 0 {
		return ""
	}
	configRepository := configrepository.Get()

	outputVariants := configRepository.GetStreamOutputVariants()

	type singleVariant struct {
		isVideoPassthrough bool
		bitrate            int
	}

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

	maxBitrate := streamSortVariants[0].bitrate
	if inboundBitrate > maxBitrate {
		return fmt.Sprintf("You're streaming to Owncast at %dkbps but only broadcasting to your viewers at %dkbps, requiring unnecessary work to be performed and possible excessive CPU use. You may want to decrease what you're sending to Owncast or increase what you send to your viewers so the highest bitrate matches.", inboundBitrate, maxBitrate)
	}

	return ""
}

func (m *Metrics) cpuUsageHealthOverviewMessage() string {
	if len(m.metrics.CPUUtilizations) < 2 {
		return ""
	}

	recentCPUUses := m.metrics.CPUUtilizations[len(m.metrics.CPUUtilizations)-2:]
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

func (m *Metrics) errorCountHealthOverviewMessage() string {
	totalNumberOfClients := len(m.windowedBandwidths)
	if totalNumberOfClients == 0 {
		return ""
	}

	clientsWithErrors := m.getClientsWithErrorsCount()

	if clientsWithErrors == 0 {
		return ""
	}

	// Only return these detailed values and messages if we feel we have enough
	// clients to be able to make a reasonable assessment. This is an arbitrary
	// number but 1 out of 1 isn't helpful.

	if totalNumberOfClients >= minClientCountForDetails {
		healthyPercentage := utils.IntPercentage(clientsWithErrors, totalNumberOfClients)

		configRepository := configrepository.Get()

		isUsingPassthrough := false
		outputVariants := configRepository.GetStreamOutputVariants()
		for _, variant := range outputVariants {
			if variant.IsVideoPassthrough {
				isUsingPassthrough = true
			}
		}

		if isUsingPassthrough {
			return fmt.Sprintf("%d of %d viewers (%d%%) are experiencing errors. You're currently using a video passthrough output, often known for causing playback issues for people. It is suggested you turn it off.", clientsWithErrors, totalNumberOfClients, healthyPercentage)
		}

		stat := status.Get()
		currentBroadcast := stat.GetCurrentBroadcast()
		if currentBroadcast != nil && currentBroadcast.LatencyLevel.SecondsPerSegment < 3 {
			return fmt.Sprintf("%d of %d viewers (%d%%) may be experiencing some issues. You may want to increase your latency buffer level in your video configuration to see if it helps.", clientsWithErrors, totalNumberOfClients, healthyPercentage)
		}

		return fmt.Sprintf("%d of %d viewers (%d%%) may be experiencing some issues.", clientsWithErrors, totalNumberOfClients, healthyPercentage)
	}

	return ""
}

func (m *Metrics) getClientsWithErrorsCount() int {
	clientsWithErrors := 0
	for _, errors := range m.windowedErrorCounts {
		if errors > 0 {
			clientsWithErrors++
		}
	}
	return clientsWithErrors
}

func (m *Metrics) getClientErrorHeathyPercentage() int {
	totalNumberOfClients := len(m.windowedErrorCounts)
	if totalNumberOfClients == 0 {
		return -1
	}

	clientsWithErrors := m.getClientsWithErrorsCount()

	if clientsWithErrors == 0 {
		return 100
	}

	pct := 100 - utils.IntPercentage(clientsWithErrors, totalNumberOfClients)

	return pct
}
