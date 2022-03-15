package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func setupPrometheusCollectors() {
	// Setup the Prometheus collectors.
	activeViewerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_viewer_count",
		Help:        "The number of viewers.",
		ConstLabels: labels,
	})

	activeChatClientCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_chat_client_count",
		Help:        "The number of connected chat clients.",
		ConstLabels: labels,
	})

	chatUserCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_total_chat_users",
		Help:        "The total number of chat users on this Owncast instance.",
		ConstLabels: labels,
	})

	currentChatMessageCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_current_chat_message_count",
		Help:        "The number of chat messages currently saved before cleanup.",
		ConstLabels: labels,
	})

	playbackErrorCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_playback_error_count",
		Help:        "Errors collected from players within this window",
		ConstLabels: labels,
	})

	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_cpu_usage",
		Help:        "CPU usage as seen internally to Owncast.",
		ConstLabels: labels,
	})
}
