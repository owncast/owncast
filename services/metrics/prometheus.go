package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (m *Metrics) setupPrometheusCollectors() {
	// Setup the Prometheus collectors.
	m.activeViewerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_viewer_count",
		Help:        "The number of viewers.",
		ConstLabels: m.labels,
	})

	m.activeChatClientCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_chat_client_count",
		Help:        "The number of connected chat clients.",
		ConstLabels: m.labels,
	})

	m.chatUserCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_total_chat_users",
		Help:        "The total number of chat users on this Owncast instance.",
		ConstLabels: m.labels,
	})

	m.currentChatMessageCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_current_chat_message_count",
		Help:        "The number of chat messages currently saved before cleanup.",
		ConstLabels: m.labels,
	})

	m.playbackErrorCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_playback_error_count",
		Help:        "Errors collected from players within this window",
		ConstLabels: m.labels,
	})

	m.cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_cpu_usage",
		Help:        "CPU usage as seen internally to Owncast.",
		ConstLabels: m.labels,
	})
}
