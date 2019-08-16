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

	// RegisterHTTPRequestDurationValues allows for pre-registering values of id,
	// duration, method and code so that the coresponding metrics are emited
	// right from the start of the recorder and not just with the first request
	RegisterHTTPRequestDurationValues(id string, method, code string)

	// RegisterHTTPResponseSizeValues allows for pre-registering values of id,
	// duration, method and code so that the coresponding metrics are emited
	// right from the start of the recorder and not just with the first request
	RegisterHTTPResponseSizeValues(id string, method, code string)

	// RegisterInflightRequestsValues allows for pre-registering values of id,
	// duration, method and code so that the coresponding metrics are emited
	// right from the start of the recorder and not just with the first request
	RegisterInflightRequestsValues(id string)
}

// Dummy is a dummy recorder.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveHTTPRequestDuration(ctx context.Context, id string, duration time.Duration, method, code string) {
}
func (dummy) ObserveHTTPResponseSize(ctx context.Context, id string, sizeBytes int64, method, code string) {
}
func (dummy) AddInflightRequests(ctx context.Context, id string, quantity int) {}

func (dummy) RegisterHTTPRequestDurationValues(id string, method, code string) {
}
func (dummy) RegisterHTTPResponseSizeValues(id string, method, code string) {
}
func (dummy) RegisterInflightRequestsValues(id string) {
}
