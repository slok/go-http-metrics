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
}

// Dummy is a dummy recorder.
var Dummy = &dummy{}

type dummy struct{}

func (dummy) ObserveHTTPRequestDuration(id string, duration time.Duration, method, code string) {}
