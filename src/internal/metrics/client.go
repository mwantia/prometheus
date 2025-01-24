package metrics

import "github.com/prometheus/client_golang/prometheus"

var MetricsClientSubsystem = "client"

var (
	ClientGeneratePromptTasksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsClientSubsystem,
			Name:      "generate_prompt_tasks_total",
			Help:      "Count the amount of generate prompt tasks",
		},
		[]string{"endpoint", "model", "style"},
	)
	ClientGeneratePromptTasksDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsClientSubsystem,
			Name:      "generate_prompt_tasks_duration_seconds",
			Help:      "Histogram of the time (in seconds) each generate prompt task took",
			Buckets:   []float64{.001, .002, .005, .01, .02, .05, .1, .2, .5, 1},
		},
		[]string{"endpoint", "model", "style"},
	)
)
