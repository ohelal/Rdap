package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Requests    *prometheus.CounterVec
	Latency     *prometheus.HistogramVec
	CacheHits   *prometheus.CounterVec
	CacheMisses *prometheus.CounterVec
	KafkaErrors prometheus.Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{
		Requests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rdap_requests_total",
				Help: "Total number of RDAP requests",
			},
			[]string{"type", "status"},
		),
		Latency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "rdap_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"type"},
		),
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rdap_cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"type"},
		),
		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rdap_cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"type"},
		),
		KafkaErrors: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "rdap_kafka_errors_total",
				Help: "Total number of Kafka errors",
			},
		),
	}

	return m
}
