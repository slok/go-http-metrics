package middleware_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	mockmiddleware "github.com/slok/go-http-metrics/internal/mocks/middleware"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
)

func TestMiddlewareMeasure(t *testing.T) {
	tests := []struct {
		name      string
		handlerID string
		config    middleware.Config
		recorder  func() metrics.Recorder
		setup     func() (metrics.Recorder, middleware.Reporter, func(t *testing.T))
	}{
		{
			name:      "Having default config with service, it should measure the metrics.",
			handlerID: "test01",
			config: middleware.Config{
				Service: "svc1",
			},
			setup: func() (metrics.Recorder, middleware.Reporter, func(t *testing.T)) {
				mrec := &mockmetrics.Recorder{}
				mrep := &mockmiddleware.Reporter{}

				// Reporter mocks.
				mrep.On("Context").Once().Return(context.TODO())
				mrep.On("StatusCode").Once().Return(418)
				mrep.On("Method").Once().Return("PATCH")
				mrep.On("BytesWritten").Once().Return(int64(42))

				// Recorder mocks.
				expProps := metrics.HTTPProperties{Service: "svc1", ID: "test01"}
				expRepProps := metrics.HTTPReqProperties{Service: "svc1", ID: "test01", Method: "PATCH", Code: "418"}

				mrec.On("AddInflightRequests", mock.Anything, expProps, 1).Once()
				mrec.On("AddInflightRequests", mock.Anything, expProps, -1).Once()
				mrec.On("ObserveHTTPRequestDuration", mock.Anything, expRepProps, mock.Anything).Once()
				mrec.On("ObserveHTTPResponseSize", mock.Anything, expRepProps, int64(42)).Once()

				return mrec, mrep, func(t *testing.T) {
					mrec.AssertExpectations(t)
					mrep.AssertExpectations(t)
				}
			},
		},
		{
			name:      "Custom labels should work",
			handlerID: "test01",
			config: middleware.Config{
				Service: "svc1",
			},
			setup: func() (metrics.Recorder, middleware.Reporter, func(t *testing.T)) {
				mrec := &mockmetrics.Recorder{}
				mrep := &mockmiddleware.CustomLabelReporter{}

				mrep.On("Context").Once().Return(context.TODO())
				mrep.On("StatusCode").Once().Return(418)
				mrep.On("Method").Once().Return("PATCH")
				mrep.On("BytesWritten").Once().Return(int64(42))
				mrep.On("CustomLabels").Once().Return([]string{"user_VIP"})

				expProps := metrics.HTTPProperties{
					Service:      "svc1",
					ID:           "test01",
					CustomLabels: []string{"user_VIP"},
				}
				expRepProps := metrics.HTTPReqProperties{
					Service:      "svc1",
					ID:           "test01",
					Method:       "PATCH",
					Code:         "418",
					CustomLabels: []string{"user_VIP"},
				}

				mrec.On("AddInflightRequests", mock.Anything, expProps, 1).Once()
				mrec.On("AddInflightRequests", mock.Anything, expProps, -1).Once()
				mrec.On("ObserveHTTPRequestDuration", mock.Anything, expRepProps, mock.Anything).Once()
				mrec.On("ObserveHTTPResponseSize", mock.Anything, expRepProps, int64(42)).Once()

				return mrec, mrep, func(t *testing.T) {
					mrec.AssertExpectations(t)
					mrep.AssertExpectations(t)
				}
			},
		},
		{
			name:      "Without having handler ID, it should measure the metrics using the request path.",
			handlerID: "",
			config:    middleware.Config{},
			setup: func() (metrics.Recorder, middleware.Reporter, func(t *testing.T)) {
				mrec := &mockmetrics.Recorder{}
				mrep := &mockmiddleware.Reporter{}

				// Reporter mocks.
				mrep.On("URLPath").Once().Return("/test/01")
				mrep.On("Context").Once().Return(context.TODO())
				mrep.On("StatusCode").Once().Return(418)
				mrep.On("Method").Once().Return("PATCH")
				mrep.On("BytesWritten").Once().Return(int64(42))

				// Recorder mocks.
				expRepProps := metrics.HTTPReqProperties{ID: "/test/01", Method: "PATCH", Code: "418"}

				mrec.On("AddInflightRequests", mock.Anything, mock.Anything, mock.Anything).Once()
				mrec.On("AddInflightRequests", mock.Anything, mock.Anything, mock.Anything).Once()
				mrec.On("ObserveHTTPRequestDuration", mock.Anything, expRepProps, mock.Anything).Once()
				mrec.On("ObserveHTTPResponseSize", mock.Anything, expRepProps, mock.Anything).Once()

				return mrec, mrep, func(t *testing.T) {
					mrec.AssertExpectations(t)
					mrep.AssertExpectations(t)
				}
			},
		},
		{
			name:      "Having grouped status code, it should measure the metrics using grouped status codes.",
			handlerID: "test01",
			config: middleware.Config{
				GroupedStatus: true,
			},
			setup: func() (metrics.Recorder, middleware.Reporter, func(t *testing.T)) {
				mrec := &mockmetrics.Recorder{}
				mrep := &mockmiddleware.Reporter{}

				// Reporter mocks.
				mrep.On("Context").Once().Return(context.TODO())
				mrep.On("StatusCode").Once().Return(418)
				mrep.On("Method").Once().Return("PATCH")
				mrep.On("BytesWritten").Once().Return(int64(42))

				// Recorder mocks.
				expRepProps := metrics.HTTPReqProperties{ID: "test01", Method: "PATCH", Code: "4xx"}

				mrec.On("AddInflightRequests", mock.Anything, mock.Anything, mock.Anything).Once()
				mrec.On("AddInflightRequests", mock.Anything, mock.Anything, mock.Anything).Once()
				mrec.On("ObserveHTTPRequestDuration", mock.Anything, expRepProps, mock.Anything).Once()
				mrec.On("ObserveHTTPResponseSize", mock.Anything, expRepProps, mock.Anything).Once()

				return mrec, mrep, func(t *testing.T) {
					mrec.AssertExpectations(t)
					mrep.AssertExpectations(t)
				}
			},
		},
		{
			name:      "Disabling inflight requests measuring, it shouldn't measure inflight metrics.",
			handlerID: "test01",
			config: middleware.Config{
				DisableMeasureInflight: true,
			},
			setup: func() (metrics.Recorder, middleware.Reporter, func(t *testing.T)) {
				mrec := &mockmetrics.Recorder{}
				mrep := &mockmiddleware.Reporter{}

				// Reporter mocks.
				mrep.On("Context").Once().Return(context.TODO())
				mrep.On("StatusCode").Once().Return(418)
				mrep.On("Method").Once().Return("PATCH")
				mrep.On("BytesWritten").Once().Return(int64(42))

				// Recorder mocks.
				expRepProps := metrics.HTTPReqProperties{ID: "test01", Method: "PATCH", Code: "418"}

				mrec.On("ObserveHTTPRequestDuration", mock.Anything, expRepProps, mock.Anything).Once()
				mrec.On("ObserveHTTPResponseSize", mock.Anything, expRepProps, mock.Anything).Once()

				return mrec, mrep, func(t *testing.T) {
					mrec.AssertExpectations(t)
					mrep.AssertExpectations(t)
				}
			},
		},
		{
			name:      "Disabling size measuring, it shouldn't measure size metrics.",
			handlerID: "test01",
			config: middleware.Config{
				DisableMeasureSize: true,
			},
			setup: func() (metrics.Recorder, middleware.Reporter, func(t *testing.T)) {
				mrec := &mockmetrics.Recorder{}
				mrep := &mockmiddleware.Reporter{}

				// Reporter mocks.
				mrep.On("Context").Once().Return(context.TODO())
				mrep.On("StatusCode").Once().Return(418)
				mrep.On("Method").Once().Return("PATCH")

				// Recorder mocks.
				expRepProps := metrics.HTTPReqProperties{ID: "test01", Method: "PATCH", Code: "418"}

				mrec.On("AddInflightRequests", mock.Anything, mock.Anything, mock.Anything).Once()
				mrec.On("AddInflightRequests", mock.Anything, mock.Anything, mock.Anything).Once()
				mrec.On("ObserveHTTPRequestDuration", mock.Anything, expRepProps, mock.Anything).Once()

				return mrec, mrep, func(t *testing.T) {
					mrec.AssertExpectations(t)
					mrep.AssertExpectations(t)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mrec, mrep, cleanup := tc.setup()

			tc.config.Recorder = mrec
			mdlw := middleware.New(tc.config)

			calledNext := false
			mdlw.Measure(tc.handlerID, mrep, func() { calledNext = true })

			cleanup(t)
			assert.True(t, calledNext)
		})
	}
}
