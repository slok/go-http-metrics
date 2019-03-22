package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/stretchr/testify/mock"
)

func getFakeHandler(statusCode int, responseBody string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})
}

func TestMiddlewareHandler(t *testing.T) {
	tests := []struct {
		name          string
		handlerID     string
		body          string
		statusCode    int
		req           *http.Request
		config        middleware.Config
		expHandlerID  string
		expMethod     string
		expSize       int64
		expStatusCode string
	}{
		{
			name:          "A default HTTP middleware should call the recorder to measure.",
			statusCode:    http.StatusAccepted,
			body:          "Я бэтмен",
			req:           httptest.NewRequest(http.MethodGet, "/test", nil),
			expHandlerID:  "/test",
			expSize:       15,
			expMethod:     http.MethodGet,
			expStatusCode: "202",
		},
		{
			name:          "Using custom ID in the middleware should call the recorder to measure with that ID.",
			handlerID:     "customID",
			body:          "I'm Batman",
			statusCode:    http.StatusTeapot,
			req:           httptest.NewRequest(http.MethodPost, "/test", nil),
			expHandlerID:  "customID",
			expSize:       10,
			expMethod:     http.MethodPost,
			expStatusCode: "418",
		},
		{
			name:          "Using grouped status code should group the status code.",
			config:        middleware.Config{GroupedStatus: true},
			statusCode:    http.StatusGatewayTimeout,
			req:           httptest.NewRequest(http.MethodPatch, "/test", nil),
			expHandlerID:  "/test",
			expMethod:     http.MethodPatch,
			expStatusCode: "5xx",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mr := &mmetrics.Recorder{}
			mr.On("ObserveHTTPRequestDuration", test.expHandlerID, mock.Anything, test.expMethod, test.expStatusCode)
			mr.On("ObserveHTTPResponseSize", test.expHandlerID, test.expSize, test.expMethod, test.expStatusCode)

			// Make the request.
			test.config.Recorder = mr
			m := middleware.New(test.config)
			h := m.Handler(test.handlerID, getFakeHandler(test.statusCode, test.body))
			h.ServeHTTP(httptest.NewRecorder(), test.req)

			mr.AssertExpectations(t)
		})
	}
}

func BenchmarkMiddlewareHandler(b *testing.B) {
	b.StopTimer()

	benchs := []struct {
		name      string
		handlerID string
		cfg       middleware.Config
	}{
		{
			name:      "benchmark with default settings.",
			handlerID: "",
			cfg:       middleware.Config{},
		},
		{
			name:      "benchmark disabling measuring size.",
			handlerID: "",
			cfg: middleware.Config{
				DisableMeasureSize: true,
			},
		},
		{
			name: "benchmark with grouped status code.",
			cfg: middleware.Config{
				GroupedStatus: true,
			},
		},
		{
			name:      "benchmark with predefined handler ID",
			handlerID: "benchmark1",
			cfg:       middleware.Config{},
		},
	}

	for _, bench := range benchs {
		b.Run(bench.name, func(b *testing.B) {
			// Prepare.
			bench.cfg.Recorder = metrics.Dummy
			m := middleware.New(bench.cfg)
			h := m.Handler(bench.handlerID, getFakeHandler(200, ""))
			r := httptest.NewRequest("GET", "/test", nil)

			// Make the requests.
			for n := 0; n < b.N; n++ {
				b.StartTimer()
				h.ServeHTTP(httptest.NewRecorder(), r)
				b.StopTimer()
			}
		})
	}
}
