package opencensus_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ocprometheus "go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"

	"github.com/slok/go-http-metrics/metrics"
	ocmetrics "github.com/slok/go-http-metrics/metrics/opencensus"
)

func TestOpenCensusRecorder(t *testing.T) {
	tests := []struct {
		name          string
		config        ocmetrics.Config
		recordMetrics func(r metrics.Recorder)
		expMetrics    []string
	}{
		{
			name:   "Default configuration should measure with the default metric style.",
			config: ocmetrics.Config{},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), "test1", 4*time.Second, http.MethodGet, "200")
				r.ObserveHTTPRequestDuration(context.TODO(), "test1", 175*time.Millisecond, http.MethodGet, "200")
				r.ObserveHTTPRequestDuration(context.TODO(), "test2", 30*time.Millisecond, http.MethodGet, "201")
				r.ObserveHTTPRequestDuration(context.TODO(), "test3", 700*time.Millisecond, http.MethodPost, "500")
				r.ObserveHTTPResponseSize(context.TODO(), "test4", 529930, http.MethodPost, "500")
				r.ObserveHTTPResponseSize(context.TODO(), "test4", 231, http.MethodPost, "500")
				r.ObserveHTTPResponseSize(context.TODO(), "test4", 99999999, http.MethodPatch, "429")
				r.AddInflightRequests(context.TODO(), "test1", 5)
				r.AddInflightRequests(context.TODO(), "test1", -3)
				r.AddInflightRequests(context.TODO(), "test2", 9)
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.05"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.1"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="5"} 2`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="10"} 2`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="+Inf"} 2`,
				`http_request_duration_seconds_count{code="200",handler="test1",method="GET"} 2`,

				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.05"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.1"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="5"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="201",handler="test2",method="GET"} 1`,

				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.05"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.1"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.25"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="0.5"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="5"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="500",handler="test3",method="POST"} 1`,

				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="100"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="1000"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="10000"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="100000"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="1e+06"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="1e+07"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="1e+08"} 1`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="1e+09"} 1`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",le="+Inf"} 1`,
				`http_response_size_bytes_count{code="429",handler="test4",method="PATCH"} 1`,

				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="100"} 0`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="1000"} 1`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="10000"} 1`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="100000"} 1`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="1e+06"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="1e+07"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="1e+08"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="1e+09"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",le="+Inf"} 2`,
				`http_response_size_bytes_count{code="500",handler="test4",method="POST"} 2`,

				`http_requests_inflight{handler="test1"} 2`,

				`http_requests_inflight{handler="test2"} 9`,
			},
		},
		{
			name: "Using custom buckets in the configuration should measure with custom buckets.",
			config: ocmetrics.Config{
				DurationBuckets: []float64{1, 2, 10, 20, 50, 200, 500, 1000, 2000, 5000, 10000},
			},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), "test1", 75*time.Minute, http.MethodGet, "200")
				r.ObserveHTTPRequestDuration(context.TODO(), "test1", 200*time.Hour, http.MethodGet, "200")
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="1"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="2"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="10"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="20"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="50"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="200"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="500"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="1000"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="2000"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="5000"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="10000"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",le="+Inf"} 2`,
				`http_request_duration_seconds_count{code="200",handler="test1",method="GET"} 2`,
			},
		},
		{
			name: "Using a custom labels in the configuration should measure with those labels.",
			config: ocmetrics.Config{
				HandlerIDLabel:  "route_id",
				StatusCodeLabel: "status_code",
				MethodLabel:     "http_method",
			},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), "test1", 4*time.Second, http.MethodGet, "200")
				r.ObserveHTTPRequestDuration(context.TODO(), "test1", 175*time.Millisecond, http.MethodGet, "200")
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.05"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.1"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="1"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="5"} 2`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="10"} 2`,
				`http_request_duration_seconds_bucket{http_method="GET",route_id="test1",status_code="200",le="+Inf"} 2`,
				`http_request_duration_seconds_count{http_method="GET",route_id="test1",status_code="200"} 2`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Initialize opencensus.
			ocexporter, err := ocprometheus.NewExporter(ocprometheus.Options{})
			require.NoError(err)
			view.RegisterExporter(ocexporter)

			// Create our recorder.
			test.config.UnregisterViewsBeforeRegister = true
			mrecorder, err := ocmetrics.NewRecorder(test.config)
			require.NoError(err)
			test.recordMetrics(mrecorder)

			// Set reporting time and wait to compute metrics view.
			view.SetReportingPeriod(1 * time.Millisecond)
			time.Sleep(10 * time.Millisecond)

			// Get the metrics handler and serve.
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/metrics", nil)
			ocexporter.ServeHTTP(rec, req)

			resp := rec.Result()

			// Check all metrics are present.
			if assert.Equal(http.StatusOK, resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)
				for _, expMetric := range test.expMetrics {
					assert.Contains(string(body), expMetric, "metric not present on the result")
				}
			}
		})
	}
}
