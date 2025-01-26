package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var MetricsNamespace = "queueverse"

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

func Setup(pool string) (*prometheus.Registry, error) {
	reg := prometheus.NewRegistry()
	wrap := prometheus.WrapRegistererWith(prometheus.Labels{
		"pool": pool,
	}, reg)

	wrap.MustRegister(collectors.NewGoCollector())
	wrap.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	wrap.MustRegister(ActivePluginsInfo)
	wrap.MustRegister(ActiveServicesInfo)
	wrap.MustRegister(ServerHttpRequestsDurationSeconds)
	wrap.MustRegister(ServerHttpRequestsTotal)
	wrap.MustRegister(ClientGeneratePromptTasksDurationSeconds)
	wrap.MustRegister(ClientGeneratePromptTasksTotal)

	return reg, nil
}

func RegisterActivePlugin(name, version, author string) {
	ActivePluginsInfo.WithLabelValues(name, version, author).Set(1)
}

func RegisterActiveService(plugin, name, t string) {
	ActiveServicesInfo.WithLabelValues(plugin, name, t).Set(1)
}
