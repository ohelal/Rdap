package errors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	errorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rdap_errors_total",
			Help: "Total number of errors by category and severity",
		},
		[]string{"category", "severity"},
	)

	errorLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rdap_error_handling_duration_seconds",
			Help:    "Error handling latency distribution",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
		},
		[]string{"category"},
	)

	retryCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rdap_error_retries_total",
			Help: "Total number of error retries",
		},
		[]string{"category"},
	)
)

func trackErrorMetrics(err *Error) {
	errorCounter.WithLabelValues(err.Category.String(), err.Severity.String()).Inc()
} 