package prometheus

import (
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
	// Registry is the registry that ill be used by the recorder to store the metrics,
	// if the default registry is not used then it will use the default one.
	Registry prometheus.Registerer
}

func (c *Config) defaults() {
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = prometheus.DefBuckets
	}

	if c.Registry == nil {
		c.Registry = prometheus.DefaultRegisterer
	}
}

type recorder struct {
	httpRequestHistogram *prometheus.HistogramVec

	cfg Config
}

// New returns a new metrics recorder that implements the recorder
// using Prometheus as the backend.
func New(cfg Config) metrics.Recorder {
	cfg.defaults()

	r := &recorder{
		httpRequestHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The latency of the HTTP requests.",
			Buckets:   cfg.DurationBuckets,
		}, []string{"handler", "method", "code"}),

		cfg: cfg,
	}

	r.registerMetrics()

	return r
}

func (r recorder) registerMetrics() {
	r.cfg.Registry.MustRegister(r.httpRequestHistogram)
}

func (r recorder) ObserveHTTPRequestDuration(id string, duration time.Duration, method, code string) {
	r.httpRequestHistogram.WithLabelValues(id, method, code).Observe(duration.Seconds())
}
