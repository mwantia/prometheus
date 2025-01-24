package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var MetricsServerSubsystem = "server"

var (
	ServerHttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsServerSubsystem,
			Name:      "http_requests_total",
			Help:      "Count the amount of http requests",
		},
		[]string{"method", "addr", "path"},
	)
	ServerHttpRequestsDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsServerSubsystem,
			Name:      "http_requests_duration_seconds",
			Help:      "Histogram of the time (in seconds) each http request took",
			Buckets:   []float64{.001, .002, .005, .01, .02, .05, .1, .2, .5, 1},
		},
		[]string{"method", "addr", "path"},
	)
)
