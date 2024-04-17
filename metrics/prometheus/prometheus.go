package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/go-http-metrics/metrics"
)

// NativeHistogramConfig provides options for tuning native histogram exposition.
// These are 1:1 with the setting exposed by the Prometheus client library:
// https://github.com/prometheus/client_golang/blob/main/prometheus/histogram.go.
// The minimal requirement is to set BucketFactor > 1 for the native histogram
// exposition to be available.
type NativeHistogramConfig struct {
	BucketFactor     float64
	ZeroThreshold    float64
	MaxBucketNumber  uint32
	MinResetDuration time.Duration
	MaxZeroThreshold float64
}

func (n NativeHistogramConfig) mapOntoOpts(opts prometheus.HistogramOpts) prometheus.HistogramOpts {
	opts.NativeHistogramBucketFactor = n.BucketFactor
	opts.NativeHistogramZeroThreshold = n.ZeroThreshold
	opts.NativeHistogramMaxBucketNumber = n.MaxBucketNumber
	opts.NativeHistogramMaxZeroThreshold = n.MaxZeroThreshold
	opts.NativeHistogramMinResetDuration = n.MinResetDuration
	return opts
}

// Config has the dependencies and values of the recorder.
type Config struct {
	// Prefix is the prefix that will be set on the metrics, by default it will be empty.
	Prefix string
	// DurationBuckets are the buckets used by Prometheus for the HTTP request duration metrics,
	// by default uses Prometheus default buckets (from 5ms to 10s).
	DurationBuckets []float64
	// DurationNativeHistogramConfig provides configuration for exposing HTTP request duration metrics
	// as Native Histograms.
	DurationNativeHistogramConfig NativeHistogramConfig
	// SizeBuckets are the buckets used by Prometheus for the HTTP response size metrics,
	// by default uses a exponential buckets from 100B to 1GB.
	SizeBuckets []float64
	// SizeNativeHistogramConfig provides configuration for exposing HTTP response size metrics
	// as Native Histograms.
	SizeNativeHistogramConfig NativeHistogramConfig
	// Registry is the registry that will be used by the recorder to store the metrics,
	// if the default registry is not used then it will use the default one.
	Registry prometheus.Registerer
	// HandlerIDLabel is the name that will be set to the handler ID label, by default is `handler`.
	HandlerIDLabel string
	// StatusCodeLabel is the name that will be set to the status code label, by default is `code`.
	StatusCodeLabel string
	// MethodLabel is the name that will be set to the method label, by default is `method`.
	MethodLabel string
	// ServiceLabel is the name that will be set to the service label, by default is `service`.
	ServiceLabel string
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

	if c.ServiceLabel == "" {
		c.ServiceLabel = "service"
	}
}

type recorder struct {
	httpRequestDurHistogram   *prometheus.HistogramVec
	httpResponseSizeHistogram *prometheus.HistogramVec
	httpRequestsInflight      *prometheus.GaugeVec
}

// NewRecorder returns a new metrics recorder that implements the recorder
// using Prometheus as the backend.
func NewRecorder(cfg Config) metrics.Recorder {
	cfg.defaults()

	r := &recorder{
		httpRequestDurHistogram: prometheus.NewHistogramVec(
			cfg.DurationNativeHistogramConfig.mapOntoOpts(prometheus.HistogramOpts{
				Namespace: cfg.Prefix,
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "The latency of the HTTP requests.",
				Buckets:   cfg.DurationBuckets,
			}), []string{cfg.ServiceLabel, cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),

		httpResponseSizeHistogram: prometheus.NewHistogramVec(
			cfg.SizeNativeHistogramConfig.mapOntoOpts(prometheus.HistogramOpts{
				Namespace: cfg.Prefix,
				Subsystem: "http",
				Name:      "response_size_bytes",
				Help:      "The size of the HTTP responses.",
				Buckets:   cfg.SizeBuckets,
			}), []string{cfg.ServiceLabel, cfg.HandlerIDLabel, cfg.MethodLabel, cfg.StatusCodeLabel}),

		httpRequestsInflight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "requests_inflight",
			Help:      "The number of inflight requests being handled at the same time.",
		}, []string{cfg.ServiceLabel, cfg.HandlerIDLabel}),
	}

	cfg.Registry.MustRegister(
		r.httpRequestDurHistogram,
		r.httpResponseSizeHistogram,
		r.httpRequestsInflight,
	)

	return r
}

func (r recorder) ObserveHTTPRequestDuration(_ context.Context, p metrics.HTTPReqProperties, duration time.Duration) {
	r.httpRequestDurHistogram.WithLabelValues(p.Service, p.ID, p.Method, p.Code).Observe(duration.Seconds())
}

func (r recorder) ObserveHTTPResponseSize(_ context.Context, p metrics.HTTPReqProperties, sizeBytes int64) {
	r.httpResponseSizeHistogram.WithLabelValues(p.Service, p.ID, p.Method, p.Code).Observe(float64(sizeBytes))
}

func (r recorder) AddInflightRequests(_ context.Context, p metrics.HTTPProperties, quantity int) {
	r.httpRequestsInflight.WithLabelValues(p.Service, p.ID).Add(float64(quantity))
}