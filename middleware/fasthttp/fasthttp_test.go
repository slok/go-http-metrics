package fasthttp_test

import (
	"testing"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	fasthttpMiddleware "github.com/slok/go-http-metrics/middleware/fasthttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		handlerID   string
		config      middleware.Config
		req         func() *fasthttp.RequestCtx
		mock        func(m *mmetrics.Recorder)
		handler     func() fasthttp.RequestHandler
		expRespCode int
		expRespBody string
	}{
		"A default HTTP middleware should call the recorder to measure.": {
			req: func() *fasthttp.RequestCtx {
				ctx := &fasthttp.RequestCtx{
					Request:  fasthttp.Request{},
					Response: fasthttp.Response{},
				}

				ctx.Request.Header.SetMethod(fasthttp.MethodPost)
				ctx.Request.Header.SetRequestURI("/test")

				return ctx
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
			handler: func() fasthttp.RequestHandler {
				return func(c *fasthttp.RequestCtx) {
					c.SetStatusCode(202)
					c.SetBodyString("test1")
				}
			},
			expRespCode: 202,
			expRespBody: "test1",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Mocks.
			mr := &mmetrics.Recorder{}
			test.mock(mr)

			// Create our instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			rCtx := test.req()

			// Make the request.
			handler := fasthttpMiddleware.Handler(test.handlerID, mdlw, test.handler())
			handler(rCtx)

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(t, test.expRespCode, rCtx.Response.StatusCode())
			assert.Equal(t, test.expRespBody, string(rCtx.Response.Body()))
		})
	}
}
