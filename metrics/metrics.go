package metrics

import (
	"time"
)

// Recorder knows how to record and measure the metrics. This
// Interface has the required methods to be used with the HTTP
// middlewares.
type Recorder interface {
	// ObserveHTTPRequestDuration measures the duration of an HTTP request.
	ObserveHTTPRequestDuration(id string, duration time.Duration, method, code string)
	// ObserveHTTPResponseSize measures the size of an HTTP response in bytes.
	ObserveHTTPResponseSize(id string, sizeBytes int64, method, code string)
	// AddInflightRequests increments and decrements the number of inflight request being
	// processed.
	AddInflightRequests(id string, quantity int)
}

// Dummy is a dummy recorder.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveHTTPRequestDuration(id string, duration time.Duration, method, code string) {}
func (dummy) ObserveHTTPResponseSize(id string, sizeBytes int64, method, code string)           {}
func (dummy) AddInflightRequests(id string, quantity int)                                       {}
