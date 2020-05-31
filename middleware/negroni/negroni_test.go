package negroni_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/negroni"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	negronimiddleware "github.com/slok/go-http-metrics/middleware/negroni"
)

func getTestHandler(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})
}

func TestMiddleware(t *testing.T) {
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

			// Create our negroni instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			n := negroni.Classic()
			n.Use(negronimiddleware.Handler("", mdlw))
			n.UseHandler(getTestHandler(test.statusCode))

			// Make the request.
			resp := httptest.NewRecorder()
			n.ServeHTTP(resp, test.req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.statusCode, resp.Result().StatusCode)
		})
	}
}
