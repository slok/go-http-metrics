package gorestful_test

import (
	"io/ioutil"
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
		handlerID   string                     // This is the name of the handler.  If empty, the path will be used.
		route       string                     // This is the route path to use go-restful.
		config      gorestfulmiddleware.Config // This is the go-restful middleware config.
		req         func() *http.Request
		mock        func(m *mmetrics.Recorder)
		handler     func() gorestful.RouteFunction
		expRespCode int
		expRespBody string
	}{
		"A default HTTP middleware should call the recorder to measure.": {
			route:  "/test",
			config: gorestfulmiddleware.Config{},
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
		"The handler ID overrides the path.": {
			handlerID: "my-handler",
			route:     "/test/{id}",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test/1", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "my-handler",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(5)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "my-handler",
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
		"The full path is used by default.": {
			route: "/test/{id}",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test/1", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test/1",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(5)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test/1",
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
		"The route path is used when desired.": {
			route: "/test/{id}",
			config: gorestfulmiddleware.Config{
				UseRoutePath: true,
			},
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/test/1", nil)
			},
			mock: func(m *mmetrics.Recorder) {
				expHTTPReqProps := metrics.HTTPReqProperties{
					ID:      "/test/{id}",
					Service: "",
					Method:  "POST",
					Code:    "202",
				}
				m.On("ObserveHTTPRequestDuration", mock.Anything, expHTTPReqProps, mock.Anything).Once()
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(5)).Once()

				expHTTPProps := metrics.HTTPProperties{
					ID:      "/test/{id}",
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
			c.Filter(gorestfulmiddleware.HandlerWithConfig(test.handlerID, mdlw, test.config))
			ws := &gorestful.WebService{}
			ws.Produces(gorestful.MIME_JSON)
			req := test.req()
			ws.Route(ws.Method(req.Method).Path(test.route).To(test.handler()))
			c.Add(ws)

			// Make the request.
			resp := httptest.NewRecorder()
			c.ServeHTTP(resp, req)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := ioutil.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}
