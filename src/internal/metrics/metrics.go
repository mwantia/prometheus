package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var MetricsNamespace = "ollama_queue"

var (
	ActivePluginsInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: MetricsNamespace,
			Name:      "active_plugins_info",
			Help:      "Active plugins information",
		},
		[]string{"name", "version", "author"},
	)
	ActiveServicesInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: MetricsNamespace,
			Name:      "active_services_info",
			Help:      "Active services information",
		},
		[]string{"plugin", "name", "type"},
	)
)

func init() {
	prometheus.MustRegister(ActivePluginsInfo)
	prometheus.MustRegister(ActiveServicesInfo)
}

func RegisterActivePlugin(name, version, author string) {
	ActivePluginsInfo.WithLabelValues(name, version, author).Set(1)
}

func RegisterActiveService(plugin, name, t string) {
	ActiveServicesInfo.WithLabelValues(plugin, name, t).Set(1)
}
