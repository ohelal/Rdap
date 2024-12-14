package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PerformanceMetrics struct {
	RequestDuration *prometheus.HistogramVec
	CacheHitRatio   *prometheus.GaugeVec
	MemoryUsage     prometheus.Gauge
	GoroutineCount  prometheus.Gauge
}

func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "rdap_request_duration_seconds",
				Help:    "Request duration distribution",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
			},
			[]string{"endpoint", "status"},
		),
		CacheHitRatio: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "rdap_cache_hit_ratio",
				Help: "Cache hit ratio by type",
			},
			[]string{"cache_type"},
		),
		MemoryUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "rdap_memory_usage_bytes",
				Help: "Current memory usage",
			},
		),
		GoroutineCount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "rdap_goroutine_count",
				Help: "Current number of goroutines",
			},
		),
	}
} 