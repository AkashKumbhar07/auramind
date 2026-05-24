package monitoring

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	registry      *prometheus.Registry
	httpRequests  *prometheus.CounterVec
	httpDuration  *prometheus.HistogramVec
	wsConnections prometheus.Gauge
	logger        *zap.Logger
}

func New(service string, logger *zap.Logger) *Metrics {
	registry := prometheus.NewRegistry()

	m := &Metrics{
		registry: registry,
		logger:   logger,
		httpRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_http_requests_total", service),
				Help: "Total HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		httpDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    fmt.Sprintf("%s_http_duration_seconds", service),
				Help:    "HTTP request duration in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		wsConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: fmt.Sprintf("%s_ws_connections", service),
				Help: "Current WebSocket connections.",
			},
		),
	}

	registry.MustRegister(m.httpRequests, m.httpDuration, m.wsConnections)

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *Metrics) IncRequests(method, path, status string) {
	m.httpRequests.WithLabelValues(method, path, status).Inc()
}

func (m *Metrics) ObserveDuration(method, path string, seconds float64) {
	m.httpDuration.WithLabelValues(method, path).Observe(seconds)
}

func (m *Metrics) SetWSConnections(n float64) {
	m.wsConnections.Set(n)
}

func (m *Metrics) ServeHTTP(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.Handler())

	addr := fmt.Sprintf(":%d", port)
	m.logger.Info("metrics server starting", zap.String("addr", addr))
	return http.ListenAndServe(addr, mux)
}
