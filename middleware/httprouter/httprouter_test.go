package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
)

func getTestHandler(statusCode int) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(statusCode)
	})
}

func TestMiddlewareIntegration(t *testing.T) {
	tests := map[string]struct {
		handlerID     string
		statusCode    int
		req           *http.Request
		config        middleware.Config
		expHandlerID  string
		expService    string
		expMethod     string
		expStatusCode string
	}{
		"A default HTTP middleware should call the recorder to measure.": {
			statusCode:    http.StatusAccepted,
			req:           httptest.NewRequest(http.MethodPost, "/test", nil),
			expHandlerID:  "/test",
			expMethod:     http.MethodPost,
			expStatusCode: "202",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mr := &mmetrics.Recorder{}
			expHTTPReqProps := metrics.HTTPReqProperties{
				ID:      test.expHandlerID,
				Service: test.expService,
				Method:  test.expMethod,
				Code:    test.expStatusCode,
			}
			expHTTPProps := metrics.HTTPProperties{
				ID:      test.expHandlerID,
				Service: test.expService,
			}
			mr.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
			mr.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, mock.Anything).Once()
			mr.On("AddInflightRequests", mock.Anything, expHTTPProps, 1).Once()
			mr.On("AddInflightRequests", mock.Anything, expHTTPProps, -1).Once()

			// Create our instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			r := httprouter.New()
			r.POST("/test", httproutermiddleware.Measure("", getTestHandler(test.statusCode), mdlw))

			// Make the request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, test.req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.statusCode, resp.Result().StatusCode)
		})
	}
}
