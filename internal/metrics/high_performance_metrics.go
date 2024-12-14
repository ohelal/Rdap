package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
)

type MetricsCollector struct {
	requestCounter  *prometheus.CounterVec
	responseLatency *prometheus.HistogramVec
	cacheHitRatio   prometheus.Gauge
	errorRate       prometheus.Gauge
	memStats        *runtime.MemStats
}

func NewMetricsCollector() *MetricsCollector {
	collector := &MetricsCollector{
		requestCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rdap_requests_total",
				Help: "Total number of RDAP requests",
			},
			[]string{"type", "status"},
		),
		responseLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "rdap_response_latency_seconds",
				Help:    "Response latency distribution",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
			},
			[]string{"type"},
		),
		cacheHitRatio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "rdap_cache_hit_ratio",
				Help: "Cache hit ratio",
			},
		),
		errorRate: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "rdap_error_rate",
				Help: "Error rate per second",
			},
		),
		memStats: &runtime.MemStats{},
	}

	// Register metrics with Prometheus
	prometheus.MustRegister(collector.requestCounter)
	prometheus.MustRegister(collector.responseLatency)
	prometheus.MustRegister(collector.cacheHitRatio)
	prometheus.MustRegister(collector.errorRate)

	return collector
}

func (m *MetricsCollector) CollectMetrics() {
	// Update memory stats
	runtime.ReadMemStats(m.memStats)
}

func (m *MetricsCollector) Close() {
	// Unregister metrics
	prometheus.Unregister(m.requestCounter)
	prometheus.Unregister(m.responseLatency)
	prometheus.Unregister(m.cacheHitRatio)
	prometheus.Unregister(m.errorRate)
}
