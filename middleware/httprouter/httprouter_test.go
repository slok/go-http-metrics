package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/middleware"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
)

func getTestHandler(statusCode int) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(statusCode)
	})
}

func TestMiddlewareIntegration(t *testing.T) {
	tests := []struct {
		name          string
		handlerID     string
		statusCode    int
		req           *http.Request
		config        middleware.Config
		expHandlerID  string
		expMethod     string
		expStatusCode string
	}{
		{
			name:          "A default HTTP middleware should call the recorder to measure.",
			statusCode:    http.StatusAccepted,
			req:           httptest.NewRequest(http.MethodPost, "/test", nil),
			expHandlerID:  "/test",
			expMethod:     http.MethodPost,
			expStatusCode: "202",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mr := &mmetrics.Recorder{}
			mr.On("ObserveHTTPRequestDuration", test.expHandlerID, mock.Anything, test.expMethod, test.expStatusCode).Once()
			mr.On("ObserveHTTPResponseSize", test.expHandlerID, mock.Anything, test.expMethod, test.expStatusCode).Once()
			mr.On("AddInflightRequests", test.expHandlerID, 1).Once()
			mr.On("AddInflightRequests", test.expHandlerID, -1).Once()

			// Create our instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			r := httprouter.New()
			r.POST("/test", httproutermiddleware.Handler("", getTestHandler(test.statusCode), mdlw))

			// Make the request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, test.req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.statusCode, resp.Result().StatusCode)
		})
	}
}
