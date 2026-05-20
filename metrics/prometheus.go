package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (s *Service) setupPrometheusCollectors() {
	// Setup the Prometheus collectors.
	s.activeViewerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_viewer_count",
		Help:        "The number of viewers.",
		ConstLabels: s.labels,
	})

	s.activeChatClientCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_active_chat_client_count",
		Help:        "The number of connected chat clients.",
		ConstLabels: s.labels,
	})

	s.chatUserCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_total_chat_users",
		Help:        "The total number of chat users on this Owncast instance.",
		ConstLabels: s.labels,
	})

	s.currentChatMessageCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_current_chat_message_count",
		Help:        "The number of chat messages currently saved before cleanup.",
		ConstLabels: s.labels,
	})

	s.playbackErrorCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_playback_error_count",
		Help:        "Errors collected from players within this window",
		ConstLabels: s.labels,
	})

	s.cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "owncast_instance_cpu_usage",
		Help:        "CPU usage as seen internally to Owncast.",
		ConstLabels: s.labels,
	})
}
