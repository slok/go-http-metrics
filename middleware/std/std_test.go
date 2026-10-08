package std_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mmetrics "github.com/slok/go-http-metrics/internal/mocks/metrics"
	"github.com/slok/go-http-metrics/metrics"
	"github.com/slok/go-http-metrics/middleware"
	stdmiddleware "github.com/slok/go-http-metrics/middleware/std"
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
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(15)).Once()

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
					w.Write([]byte("Я бэтмен")) // nolint: errcheck
				})
			},
			expRespCode: 202,
			expRespBody: "Я бэтмен",
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
			test.config.Recorder = mr
			m := middleware.New(test.config)
			h := stdmiddleware.Handler(test.handlerID, m, test.handler())

			// Make the request.
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, test.req())

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := io.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}

func TestProvider(t *testing.T) {
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
				m.On("ObserveHTTPResponseSize", mock.Anything, expHTTPReqProps, int64(15)).Once()

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
					w.Write([]byte("Я бэтмен")) // nolint: errcheck
				})
			},
			expRespCode: 202,
			expRespBody: "Я бэтмен",
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
			test.config.Recorder = mr
			m := middleware.New(test.config)
			provider := stdmiddleware.HandlerProvider(test.handlerID, m)
			h := provider(test.handler())

			// Make the request.
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, test.req())

			// Check.
			mr.AssertExpectations(t)
			assert.Equal(test.expRespCode, resp.Result().StatusCode)
			gotBody, err := io.ReadAll(resp.Result().Body)
			require.NoError(err)
			assert.Equal(test.expRespBody, string(gotBody))
		})
	}
}

type writeResult struct {
	n   int
	err error
}

type responseSizeWriter struct {
	*httptest.ResponseRecorder
	results []writeResult
	writes  int
}

func (w *responseSizeWriter) Write(p []byte) (int, error) {
	result := w.results[w.writes]
	w.writes++
	_, _ = w.ResponseRecorder.Write(p[:result.n])
	return result.n, result.err
}

func TestResponseSize(t *testing.T) {
	tests := map[string]struct {
		results []writeResult
		expSize int64
	}{
		"Full writes": {
			results: []writeResult{{n: 5}, {n: 5}},
			expSize: 10,
		},
		"Partial write with an error": {
			results: []writeResult{{n: 2, err: io.ErrShortWrite}},
			expSize: 2,
		},
		"Rejected write": {
			results: []writeResult{{err: http.ErrBodyNotAllowed}},
			expSize: 0,
		},
		"Full write followed by a partial write": {
			results: []writeResult{{n: 5}, {n: 3, err: io.ErrShortWrite}},
			expSize: 8,
		},
		"Full write with an error": {
			results: []writeResult{{n: 5, err: io.ErrClosedPipe}},
			expSize: 5,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mr := &mmetrics.Recorder{}
			mr.On("ObserveHTTPRequestDuration", mock.Anything, mock.Anything, mock.Anything).Once()
			mr.On("ObserveHTTPResponseSize", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				assert.Equal(t, test.expSize, args.Get(2))
			}).Once()
			mr.On("AddInflightRequests", mock.Anything, mock.Anything, 1).Once()
			mr.On("AddInflightRequests", mock.Anything, mock.Anything, -1).Once()

			m := middleware.New(middleware.Config{Recorder: mr})
			h := stdmiddleware.Handler("", m, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, result := range test.results {
					n, err := w.Write([]byte("hello"))
					assert.Equal(t, result.n, n)
					assert.Equal(t, result.err, err)
				}
			}))
			w := &responseSizeWriter{
				ResponseRecorder: httptest.NewRecorder(),
				results:          test.results,
			}
			h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

			assert.Equal(t, int(test.expSize), w.Body.Len())
			mr.AssertExpectations(t)
		})
	}
}
