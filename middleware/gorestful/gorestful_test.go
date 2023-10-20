package gorestful_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gorestful "github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	gorestfulmiddleware "github.com/slok/go-http-metrics/middleware/gorestful"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		handlerID   string
		config      middleware.Config
		req         func() *http.Request
		mock        func(m *mmetrics.Recorder)
		handler     func() gorestful.RouteFunction
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
			handler: func() gorestful.RouteFunction {
				return gorestful.RouteFunction(func(_ *gorestful.Request, resp *gorestful.Response) {
					resp.WriteHeader(202)
					resp.Write([]byte("test1")) // nolint: errcheck
				})
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
			c := gorestful.NewContainer()
			c.Filter(gorestfulmiddleware.Handler(test.handlerID, mdlw))
			ws := &gorestful.WebService{}
			ws.Produces(gorestful.MIME_JSON)
			req := test.req()
			ws.Route(ws.Method(req.Method).Path(req.URL.Path).To(test.handler()))
			c.Add(ws)

			// Make the request.
			resp := httptest.NewRecorder()
			c.ServeHTTP(resp, req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := io.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}
