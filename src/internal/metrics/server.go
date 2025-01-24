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
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 15),
		},
		[]string{"method", "addr", "path"},
	)
)
