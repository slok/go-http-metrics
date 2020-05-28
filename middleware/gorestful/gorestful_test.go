package gorestful_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gorestful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	gorestfulmiddleware "github.com/slok/go-http-metrics/middleware/gorestful"
)

func getTestHandler(statusCode int) gorestful.RouteFunction {
	return gorestful.RouteFunction(func(_ *gorestful.Request, resp *gorestful.Response) {
		_ = resp.WriteHeaderAndEntity(statusCode, "Hello world")
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
			c := gorestful.NewContainer()
			c.Filter(gorestfulmiddleware.Measure("", mdlw))
			ws := &gorestful.WebService{}
			ws.Produces(gorestful.MIME_JSON)
			ws.Route(ws.POST("/test").To(getTestHandler(test.statusCode)))
			c.Add(ws)

			// Make the request.
			resp := httptest.NewRecorder()
			c.ServeHTTP(resp, test.req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.statusCode, resp.Result().StatusCode)
		})
	}
}
