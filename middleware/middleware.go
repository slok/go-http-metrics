// Package middleware will measure metrics of a Go net/http
// handler in Prometheus format.
// The metrics measured are based on RED and/or Four golden signals and
// try to be measured in a efficient way.
package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/slok/go-http-metrics/metrics"
)

// Config is the configuration for the middleware factory.
type Config struct {
	// Recorder is the way the metrics will be recorder in the different backends.
	Recorder metrics.Recorder
	// GroupedStatus will group the status label in the form of `\dxx`, for example,
	// 200, 201, and 203 will have the label `code="2xx"`. This impacts on the cardinality
	// of the metrics and also improves the performance of queries that are grouped by
	// status code because there are already aggregated in the metric.
	// By default will be false.
	GroupedStatus bool
}

func (c *Config) validate() {
	if c.Recorder == nil {
		c.Recorder = metrics.Dummy
	}
}

// Middleware is a factory that creates middlewares or wrappers that
// measure requests to the wrapped handler using different metrics
// backends using a `metrics.Recorder` implementation.
type Middleware interface {
	// Handler wraps the received handler with the Prometheus middleware.
	// The first argument receives the handlerID, all the metrics will have
	// that handler ID as the handler label on the metrics, if an empty
	// string is passed then it will get the handlerID from the request
	// path.
	Handler(handlerID string, h http.Handler) http.Handler
}

// middelware is the prometheus middleware instance.
type middleware struct {
	cfg Config
}

// New returns the a Middleware factory.
func New(cfg Config) Middleware {
	// Validate the configuration.
	cfg.validate()

	// Create our middleware with all the configuration options.
	m := &middleware{
		cfg: cfg,
	}

	return m
}

// Handler satisfies Middleware interface.
func (m *middleware) Handler(handlerID string, h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intercept the writer so we can retrieve data afterwards.
		wi := &responseWriterInterceptor{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		// If there isn't predefined handler ID we
		// set that ID as the URL path.
		hid := handlerID
		if handlerID == "" {
			hid = r.URL.Path
		}

		// Start the timer and when finishing measure the duration.
		start := time.Now()
		defer func() {
			duration := time.Since(start)

			// If we need to group the status code, it uses the
			// first number of the status code because is the least
			// required identification way.
			var code string
			if m.cfg.GroupedStatus {
				code = fmt.Sprintf("%dxx", wi.statusCode/100)
			} else {
				code = strconv.Itoa(wi.statusCode)
			}

			m.cfg.Recorder.ObserveHTTPRequestDuration(hid, duration, r.Method, code)
		}()

		h.ServeHTTP(wi, r)
	})
}

// responseWriterInterceptor is a simple wrapper to intercept set data on a
// ResponseWriter.
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
