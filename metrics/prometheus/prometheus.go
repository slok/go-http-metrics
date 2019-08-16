package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/go-http-metrics/metrics"
)

// Config has the dependencies and values of the recorder.
type Config struct {
	// Prefix is the prefix that will be set on the metrics, by default it will be empty.
	Prefix string
	// DurationBuckets are the buckets used by Prometheus for the HTTP request duration metrics,
	// by default uses Prometheus default buckets (from 5ms to 10s).
	DurationBuckets []float64
	// SizeBuckets are the buckets used by Prometheus for the HTTP response size metrics,
	// by default uses a exponential buckets from 100B to 1GB.
	SizeBuckets []float64
	// Registry is the registry that will be used by the recorder to store the metrics,
	// if the default registry is not used then it will use the default one.
	Registry prometheus.Registerer
	// HandlerIDLabel is the name that will be set to the handler ID label, by default is `handler`.
	HandlerIDLabel string
	// StatusCodeLabel is the name that will be set to the status code label, by default is `code`.
	StatusCodeLabel string
	// MethodLabel is the name that will be set to the method label, by default is `method`.
	MethodLabel string
}

func (c *Config) defaults() {
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = prometheus.DefBuckets
	}

	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = prometheus.ExponentialBuckets(100, 10, 8)
	}

	if c.Registry == nil {
		c.Registry = prometheus.DefaultRegisterer
	}

	if c.HandlerIDLabel == "" {
		c.HandlerIDLabel = "handler"
	}

	if c.StatusCodeLabel == "" {
		c.StatusCodeLabel = "code"
	}

	if c.MethodLabel == "" {
		c.MethodLabel = "method"
	}
}

type recorder struct {
	httpRequestDurHistogram   *prometheus.HistogramVec
	httpResponseSizeHistogram *prometheus.HistogramVec
	httpRequestsInflight      *prometheus.GaugeVec

	cfg Config
}

// NewRecorder returns a new metrics recorder that implements the recorder
// using Prometheus as the backend.
func NewRecorder(cfg Config) metrics.Recorder {
	cfg.defaults()

	r := &recorder{
		httpRequestDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The latency of the HTTP requests.",
			Buckets:   cfg.DurationBuckets,
		}, []string{cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),
		httpResponseSizeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "The size of the HTTP responses.",
			Buckets:   cfg.SizeBuckets,
		}, []string{cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),
		httpRequestsInflight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "requests_inflight",
			Help:      "The number of inflight requests being handled at the same time.",
		}, []string{cfg.HandlerIDLabel}),

		cfg: cfg,
	}

	r.registerMetrics()

	return r
}

func (r recorder) registerMetrics() {
	r.cfg.Registry.MustRegister(
		r.httpRequestDurHistogram,
		r.httpResponseSizeHistogram,
		r.httpRequestsInflight,
	)
}

func (r recorder) ObserveHTTPRequestDuration(_ context.Context, id string, duration time.Duration, method, code string) {
	r.httpRequestDurHistogram.WithLabelValues(id, method, code).Observe(duration.Seconds())
}

func (r recorder) ObserveHTTPResponseSize(_ context.Context, id string, sizeBytes int64, method, code string) {
	r.httpResponseSizeHistogram.WithLabelValues(id, method, code).Observe(float64(sizeBytes))
}

func (r recorder) AddInflightRequests(_ context.Context, id string, quantity int) {
	r.httpRequestsInflight.WithLabelValues(id).Add(float64(quantity))
}

// RegisterHTTPRequestDurationValues allows for pre-registering values of id,
// duration, method and code so that the coresponding metrics are emited
// right from the start of the recorder and not just with the first request
func (r recorder) RegisterHTTPRequestDurationValues(id string, method, code string) {
	r.httpRequestDurHistogram.WithLabelValues(id, method, code)
}

// RegisterHTTPResponseSizeValues allows for pre-registering values of id,
// duration, method and code so that the coresponding metrics are emited
// right from the start of the recorder and not just with the first request
func (r recorder) RegisterHTTPResponseSizeValues(id string, method, code string) {
	r.httpResponseSizeHistogram.WithLabelValues(id, method, code)
}

// RegisterInflightRequestsValues allows for pre-registering values of id,
// duration, method and code so that the coresponding metrics are emited
// right from the start of the recorder and not just with the first request
func (r recorder) RegisterInflightRequestsValues(id string) {
	r.httpRequestsInflight.WithLabelValues(id)
}
