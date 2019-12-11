package opencensus

import (
	"context"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/slok/go-http-metrics/metrics"
)

var (
	durationBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	sizeBuckets     = []float64{100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
)

// Config has the dependencies and values of the recorder.
type Config struct {
	// DurationBuckets are the buckets used for the HTTP request duration metrics,
	// by default uses default buckets (from 5ms to 10s).
	DurationBuckets []float64
	// SizeBuckets are the buckets for the HTTP response size metrics,
	// by default uses a exponential buckets from 100B to 1GB.
	SizeBuckets []float64
	// HandlerIDLabel is the name that will be set to the handler ID label, by default is `handler`.
	HandlerIDLabel string
	// StatusCodeLabel is the name that will be set to the status code label, by default is `code`.
	StatusCodeLabel string
	// MethodLabel is the name that will be set to the method label, by default is `method`.
	MethodLabel string
	// ServiceLabel is the name that will be set to the service label, by default is `service`.
	ServiceLabel string
	// UnregisterViewsBeforeRegister will unregister the previous Recorder views before registering
	// again. This is required on cases where multiple instances of recorder will be made due to how
	// Opencensus is implemented (everything is at global state). Sadly this option is a kind of hack
	// so we can test without exposing the views to the user. On regular usage this option is very
	// rare to use it.
	UnregisterViewsBeforeRegister bool
}

func (c *Config) defaults() {
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = durationBuckets
	}

	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = sizeBuckets
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
	// Keys.
	codeKey    tag.Key
	methodKey  tag.Key
	handlerKey tag.Key
	serviceKey tag.Key

	// Measures.
	latencySecs   *stats.Float64Measure
	sizeBytes     *stats.Int64Measure
	inflightCount *stats.Int64Measure
}

// NewRecorder returns a new Recorder that uses OpenCensus stats
// as the backend.
func NewRecorder(cfg Config) (metrics.Recorder, error) {
	cfg.defaults()

	r := &recorder{}

	// Prepare metrics.
	err := r.createKeys(cfg)
	if err != nil {
		return nil, err
	}
	r.createMeasurements()
	err = r.registerViews(cfg)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *recorder) createKeys(cfg Config) error {
	code, err := tag.NewKey(cfg.StatusCodeLabel)
	if err != nil {
		return err
	}
	r.codeKey = code

	method, err := tag.NewKey(cfg.MethodLabel)
	if err != nil {
		return err
	}
	r.methodKey = method

	handler, err := tag.NewKey(cfg.HandlerIDLabel)
	if err != nil {
		return err
	}
	r.handlerKey = handler

	service, err := tag.NewKey(cfg.ServiceLabel)
	if err != nil {
		return err
	}
	r.serviceKey = service

	return nil
}

func (r *recorder) createMeasurements() {
	r.latencySecs = stats.Float64(
		"http_request_duration_seconds",
		"The latency of the HTTP requests",
		"s")
	r.sizeBytes = stats.Int64(
		"http_response_size_bytes",
		"The size of the HTTP responses",
		stats.UnitBytes)
	r.inflightCount = stats.Int64(
		"http_requests_inflight",
		"The number of inflight requests being handled at the same time",
		stats.UnitNone)
}

func (r recorder) registerViews(cfg Config) error {

	// OpenCensus uses global states, sadly we can't have view instance.
	durationView := &view.View{
		Name:        "http_request_duration_seconds",
		Description: "The latency of the HTTP requests",
		TagKeys:     []tag.Key{r.serviceKey, r.handlerKey, r.methodKey, r.codeKey},
		Measure:     r.latencySecs,
		Aggregation: view.Distribution(cfg.DurationBuckets...),
	}
	sizeView := &view.View{
		Name:        "http_response_size_bytes",
		Description: "The size of the HTTP responses",
		TagKeys:     []tag.Key{r.serviceKey, r.handlerKey, r.methodKey, r.codeKey},
		Measure:     r.sizeBytes,
		Aggregation: view.Distribution(cfg.SizeBuckets...),
	}
	inflightView := &view.View{
		Name:        "http_requests_inflight",
		Description: "The number of inflight requests being handled at the same time",
		TagKeys:     []tag.Key{r.serviceKey, r.handlerKey},
		Measure:     r.inflightCount,
		Aggregation: view.Sum(),
	}

	// Do we need to unregister the same views before registering.
	if cfg.UnregisterViewsBeforeRegister {
		view.Unregister(durationView, sizeView, inflightView)
	}

	err := view.Register(durationView, sizeView, inflightView)
	if err != nil {
		return err
	}

	return nil
}

func (r recorder) ObserveHTTPRequestDuration(ctx context.Context, p metrics.HTTPReqProperties, duration time.Duration) {
	ctx = r.ctxWithTagFromHTTPReqProperties(ctx, p)
	stats.Record(ctx, r.latencySecs.M(duration.Seconds()))
}

func (r recorder) ObserveHTTPResponseSize(ctx context.Context, p metrics.HTTPReqProperties, sizeBytes int64) {
	ctx = r.ctxWithTagFromHTTPReqProperties(ctx, p)
	stats.Record(ctx, r.sizeBytes.M(sizeBytes))
}

func (r recorder) AddInflightRequests(ctx context.Context, p metrics.HTTPProperties, quantity int) {
	ctx = r.ctxWithTagFromHTTPProperties(ctx, p)
	stats.Record(ctx, r.inflightCount.M(int64(quantity)))
}

func (r recorder) ctxWithTagFromHTTPReqProperties(ctx context.Context, p metrics.HTTPReqProperties) context.Context {
	newCtx, _ := tag.New(ctx,
		tag.Upsert(r.serviceKey, p.Service),
		tag.Upsert(r.handlerKey, p.ID),
		tag.Upsert(r.methodKey, p.Method),
		tag.Upsert(r.codeKey, p.Code),
	)
	return newCtx
}

func (r recorder) ctxWithTagFromHTTPProperties(ctx context.Context, p metrics.HTTPProperties) context.Context {
	newCtx, _ := tag.New(ctx,
		tag.Upsert(r.serviceKey, p.Service),
		tag.Upsert(r.handlerKey, p.ID),
	)
	return newCtx
}
