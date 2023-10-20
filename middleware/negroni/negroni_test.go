package negroni_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	negronimiddleware "github.com/slok/go-http-metrics/middleware/negroni"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		handlerID   string
		config      middleware.Config
		req         func() *http.Request
		mock        func(m *mmetrics.Recorder)
		handler     func() http.Handler
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
			handler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(202)
					w.Write([]byte("test1")) // nolint: errcheck
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

			// Create our negroni instance with the middleware.
			mdlw := middleware.New(middleware.Config{Recorder: mr})
			n := negroni.Classic()
			n.Use(negronimiddleware.Handler(test.handlerID, mdlw))
			n.UseHandler(test.handler())

			// Make the request.
			resp := httptest.NewRecorder()
			n.ServeHTTP(resp, test.req())

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := io.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}
