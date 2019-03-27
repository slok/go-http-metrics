package gorestful_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gorestful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/middleware"
	gorestfulmiddleware "github.com/slok/go-http-metrics/middleware/gorestful"
)

func getTestHandler(statusCode int) gorestful.RouteFunction {
	return gorestful.RouteFunction(func(_ *gorestful.Request, resp *gorestful.Response) {
		resp.WriteHeaderAndEntity(statusCode, "Hello world")
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
			mr.On("ObserveHTTPRequestDuration", mock.Anything, test.expHandlerID, mock.Anything, test.expMethod, test.expStatusCode).Once()
			mr.On("ObserveHTTPResponseSize", mock.Anything, test.expHandlerID, mock.Anything, test.expMethod, test.expStatusCode).Once()
			mr.On("AddInflightRequests", mock.Anything, test.expHandlerID, 1).Once()
			mr.On("AddInflightRequests", mock.Anything, test.expHandlerID, -1).Once()

			// Create our instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			c := gorestful.NewContainer()
			c.Filter(gorestfulmiddleware.Handler("", mdlw))
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
