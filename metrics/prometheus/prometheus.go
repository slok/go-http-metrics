package prometheus

import (
	"context"
	"fmt"
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
	// ServiceLabel is the name that will be set to the service label, by default is `service`.
	ServiceLabel string
	// ExtraLabels is used to add extra labels to the metrics that can be extracted per request.
	ExtraLabels []string
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
	cfg                       *Config
}

// NewRecorder returns a new metrics recorder that implements the recorder
// using Prometheus as the backend.
func NewRecorder(cfg Config) metrics.Recorder {
	cfg.defaults()

	httpRequestDurHistogramLabels := []string{cfg.ServiceLabel, cfg.HandlerIDLabel,
		cfg.MethodLabel, cfg.StatusCodeLabel}
	httpResponseSizeHistogramLabels := []string{cfg.ServiceLabel, cfg.HandlerIDLabel,
		cfg.MethodLabel, cfg.StatusCodeLabel}
	httpRequestsInflightLabels := []string{cfg.ServiceLabel, cfg.HandlerIDLabel}

	httpRequestDurHistogramLabels = append(httpRequestDurHistogramLabels, cfg.ExtraLabels...)
	httpResponseSizeHistogramLabels = append(httpResponseSizeHistogramLabels, cfg.ExtraLabels...)
	httpRequestsInflightLabels = append(httpRequestsInflightLabels, cfg.ExtraLabels...)

	r := &recorder{
		httpRequestDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The latency of the HTTP requests.",
			Buckets:   cfg.DurationBuckets,
		}, httpRequestDurHistogramLabels),

		httpResponseSizeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "The size of the HTTP responses.",
			Buckets:   cfg.SizeBuckets,
		}, httpResponseSizeHistogramLabels),

		httpRequestsInflight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "requests_inflight",
			Help:      "The number of inflight requests being handled at the same time.",
		}, httpRequestsInflightLabels),

		cfg: &cfg,
	}

	cfg.Registry.MustRegister(
		r.httpRequestDurHistogram,
		r.httpResponseSizeHistogram,
		r.httpRequestsInflight,
	)

	return r
}

func (r recorder) ObserveHTTPRequestDuration(_ context.Context, p metrics.HTTPReqProperties, duration time.Duration) {
	labelValues := r.fillLabelValues(p.ExtraProperties, p.Service, p.ID, p.Method, p.Code)
	r.httpRequestDurHistogram.WithLabelValues(labelValues...).Observe(duration.Seconds())
}

func (r recorder) ObserveHTTPResponseSize(_ context.Context, p metrics.HTTPReqProperties, sizeBytes int64) {
	labelValues := r.fillLabelValues(p.ExtraProperties, p.Service, p.ID, p.Method, p.Code)
	r.httpResponseSizeHistogram.WithLabelValues(labelValues...).Observe(float64(sizeBytes))
}

func (r recorder) AddInflightRequests(_ context.Context, p metrics.HTTPProperties, quantity int) {
	labelValues := r.fillLabelValues(p.ExtraProperties, p.Service, p.ID)
	r.httpRequestsInflight.WithLabelValues(labelValues...).Add(float64(quantity))
}

func (r recorder) fillLabelValues(extraProperties map[string]interface{}, labels ...string) []string {
	for _, extraLabel := range r.cfg.ExtraLabels {
		prop := extraProperties[extraLabel]
		if prop == nil {
			prop = ""
		}
		labels = append(labels, fmt.Sprintf("%v", prop))
	}
	return labels
}
