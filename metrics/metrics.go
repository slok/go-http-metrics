package metrics

import (
	"context"
	"time"
)

// Recorder knows how to record and measure the metrics. This
// Interface has the required methods to be used with the HTTP
// middlewares.
type Recorder interface {
	// ObserveHTTPRequestDuration measures the duration of an HTTP request.
	ObserveHTTPRequestDuration(ctx context.Context, id string, duration time.Duration, method, code string)
	// ObserveHTTPResponseSize measures the size of an HTTP response in bytes.
	ObserveHTTPResponseSize(ctx context.Context, id string, sizeBytes int64, method, code string)
	// AddInflightRequests increments and decrements the number of inflight request being
	// processed.
	AddInflightRequests(ctx context.Context, id string, quantity int)
}

// Dummy is a dummy recorder.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveHTTPRequestDuration(ctx context.Context, id string, duration time.Duration, method, code string) {
}
func (dummy) ObserveHTTPResponseSize(ctx context.Context, id string, sizeBytes int64, method, code string) {
}
func (dummy) AddInflightRequests(ctx context.Context, id string, quantity int) {}
