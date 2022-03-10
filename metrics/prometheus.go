package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels                  map[string]string
	activeViewerCount       prometheus.Gauge
	activeChatClientCount   prometheus.Gauge
	cpuUsage                prometheus.Gauge
	chatUserCount           prometheus.Gauge
	currentChatMessageCount prometheus.Gauge
	playbackErrorCount      prometheus.Gauge
)
