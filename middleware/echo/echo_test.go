package echo_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	echoMiddleware "github.com/slok/go-http-metrics/middleware/echo"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		handlerID   string
		config      middleware.Config
		req         func() *http.Request
		mock        func(m *mmetrics.Recorder)
		handler     func() echo.HandlerFunc
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
			handler: func() echo.HandlerFunc {
				return func(context echo.Context) error {
					resp := context.Response()
					resp.WriteHeader(202)
					_, err := resp.Write([]byte("test1"))
					return err
				}
			},
			expRespCode: 202,
			expRespBody: "test1",
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
			e := echo.New()
			req := test.req()
			e.Add(req.Method, req.URL.Path, test.handler(), echoMiddleware.Handler(test.handlerID, mdlw))

			// Make the request.
			resp := httptest.NewRecorder()
			e.ServeHTTP(resp, test.req())

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := io.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}
