package gin_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		handlerID   string
		config      middleware.Config
		route       string
		req         func() *http.Request
		mock        func(m *mmetrics.Recorder)
		handler     func() gin.HandlerFunc
		expRespCode int
		expRespBody string
	}{
		"A default HTTP middleware should call the recorder to measure.": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(5)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test",
					Service: "",
				}
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, 1).Once()
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, -1).Once()
			},
			handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					c.String(202, "test1")
				}
			},
			expRespCode: 202,
			expRespBody: "test1",
		},

		"A default HTTP middleware using JSON should call the recorder to measure (Regression test: https://github.com/slok/go-http-metrics/issues/31).": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(14)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test",
					Service: "",
				}
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, 1).Once()
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, -1).Once()
			},
			handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					c.JSON(202, map[string]string{"test": "one"})
				}
			},
			expRespCode: 202,
			expRespBody: `{"test":"one"}`,
		},

		"A default HTTP middleware should set template route as route label": {
			route: "/test/:id",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test/12", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test/:id",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(14)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test/:id",
					Service: "",
				}
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, 1).Once()
				m.On("AddInflightRequests", mock.Anything, expHTTPProps, -1).Once()
			},
			handler: func() gin.HandlerFunc {
				return func(c *gin.Context) {
					c.JSON(202, map[string]string{"test": "one"})
				}
			},
			expRespCode: 202,
			expRespBody: `{"test":"one"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Mocks.
			mr := &mmetrics.Recorder{}
			test.mock(mr)

			// Create our instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			engine := gin.New()
			req := test.req()

			relativePath := req.URL.Path
			if test.route != "" {
				relativePath = test.route
			}
			engine.Handle(req.Method, relativePath,
				ginmiddleware.Handler(test.handlerID, mdlw),
				test.handler())

			// Make the request.
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := io.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}
