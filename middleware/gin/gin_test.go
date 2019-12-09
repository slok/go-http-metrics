package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

func getTestHandler(statusCode int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(statusCode, "Hello world")
	}
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
			engine := gin.New()
			engine.POST("/test",
				ginmiddleware.Handler("", mdlw),
				getTestHandler(test.statusCode))

			// Make the request.
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, test.req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.statusCode, resp.Result().StatusCode)
		})
	}
}
