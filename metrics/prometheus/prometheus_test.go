package prometheus_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	"github.com/slok/go-http-metrics/metrics"
	libprometheus "github.com/slok/go-http-metrics/metrics/prometheus"
)

func TestPrometheusRecorder(t *testing.T) {
	tests := []struct {
		name          string
		config        libprometheus.Config
		recordMetrics func(r metrics.Recorder)
		expMetrics    []string
	}{
		{
			name:   "Default configuration should measure with the default metric style.",
			config: libprometheus.Config{},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 5*time.Second)
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 175*time.Millisecond)
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test2", Method: http.MethodGet, Code: "201"}, 50*time.Millisecond)
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc2", ID: "test3", Method: http.MethodPost, Code: "500"}, 700*time.Millisecond)
				r.ObserveHTTPResponseSize(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test4", Method: http.MethodPost, Code: "500"}, 529930)
				r.ObserveHTTPResponseSize(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test4", Method: http.MethodPost, Code: "500"}, 231)
				r.ObserveHTTPResponseSize(context.TODO(), metrics.HTTPReqProperties{Service: "svc2", ID: "test4", Method: http.MethodPatch, Code: "429"}, 99999999)
				r.AddInflightRequests(context.TODO(), metrics.HTTPProperties{Service: "svc1", ID: "test1"}, 5)
				r.AddInflightRequests(context.TODO(), metrics.HTTPProperties{Service: "svc1", ID: "test1"}, -3)
				r.AddInflightRequests(context.TODO(), metrics.HTTPProperties{Service: "svc2", ID: "test2"}, 9)
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.05"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.1"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="5"} 2`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="10"} 2`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="+Inf"} 2`,
				`http_request_duration_seconds_count{code="200",handler="test1",method="GET",service="svc1"} 2`,

				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.05"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.1"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="5"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="201",handler="test2",method="GET",service="svc1",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="201",handler="test2",method="GET",service="svc1"} 1`,

				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.05"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.1"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.25"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="0.5"} 0`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="5"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="500",handler="test3",method="POST",service="svc2",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="500",handler="test3",method="POST",service="svc2"} 1`,

				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="100"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="1000"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="10000"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="100000"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="1e+06"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="1e+07"} 0`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="1e+08"} 1`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="1e+09"} 1`,
				`http_response_size_bytes_bucket{code="429",handler="test4",method="PATCH",service="svc2",le="+Inf"} 1`,
				`http_response_size_bytes_count{code="429",handler="test4",method="PATCH",service="svc2"} 1`,

				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="100"} 0`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="1000"} 1`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="10000"} 1`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="100000"} 1`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="1e+06"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="1e+07"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="1e+08"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="1e+09"} 2`,
				`http_response_size_bytes_bucket{code="500",handler="test4",method="POST",service="svc1",le="+Inf"} 2`,
				`http_response_size_bytes_count{code="500",handler="test4",method="POST",service="svc1"} 2`,

				`http_requests_inflight{handler="test1",service="svc1"} 2`,

				`http_requests_inflight{handler="test2",service="svc2"} 9`,
			},
		},
		{
			name: "Using a prefix in the configuration should measure with prefix.",
			config: libprometheus.Config{
				Prefix: "batman",
			},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 5*time.Second)
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 175*time.Millisecond)
			},
			expMetrics: []string{
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.005"} 0`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.01"} 0`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.025"} 0`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.05"} 0`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.1"} 0`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.25"} 1`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="0.5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="1"} 1`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="2.5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="5"} 2`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="10"} 2`,
				`batman_http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="+Inf"} 2`,
				`batman_http_request_duration_seconds_count{code="200",handler="test1",method="GET",service="svc1"} 2`,
			},
		},
		{
			name: "Using custom buckets in the configuration should measure with custom buckets.",
			config: libprometheus.Config{
				DurationBuckets: []float64{1, 2, 10, 20, 50, 200, 500, 1000, 2000, 5000, 10000},
			},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 75*time.Minute)
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 200*time.Hour)
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="1"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="2"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="10"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="20"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="50"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="200"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="500"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="1000"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="2000"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="5000"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="10000"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",le="+Inf"} 2`,
				`http_request_duration_seconds_count{code="200",handler="test1",method="GET",service="svc1"} 2`,
			},
		},
		{
			name: "Using a custom label names in the configuration should measure with those labels.",
			config: libprometheus.Config{
				HandlerIDLabel:  "route_id",
				StatusCodeLabel: "status_code",
				MethodLabel:     "http_method",
				ServiceLabel:    "http_service",
			},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 6*time.Second)
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{Service: "svc1", ID: "test1", Method: http.MethodGet, Code: "200"}, 175*time.Millisecond)
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.005"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.01"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.025"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.05"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.1"} 0`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="1"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="5"} 1`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="10"} 2`,
				`http_request_duration_seconds_bucket{http_method="GET",http_service="svc1",route_id="test1",status_code="200",le="+Inf"} 2`,
				`http_request_duration_seconds_count{http_method="GET",http_service="svc1",route_id="test1",status_code="200"} 2`,
			},
		},
		{
			name: "Using a custom labels in the configuration should measure with those labels.",
			config: libprometheus.Config{
				DurationBuckets: []float64{1, 10},
				CustomLabels:    []string{"user_id"},
			},
			recordMetrics: func(r metrics.Recorder) {
				r.ObserveHTTPRequestDuration(context.TODO(), metrics.HTTPReqProperties{
					Service:      "svc1",
					ID:           "test1",
					Method:       http.MethodGet,
					Code:         "200",
					CustomLabels: []string{"userVIP"},
				}, 6*time.Second)
				r.AddInflightRequests(context.TODO(), metrics.HTTPProperties{
					Service:      "svc1",
					ID:           "test1",
					CustomLabels: []string{"userVIP"},
				}, 1)
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",user_id="userVIP",le="1"} 0`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",user_id="userVIP",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="200",handler="test1",method="GET",service="svc1",user_id="userVIP",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="200",handler="test1",method="GET",service="svc1",user_id="userVIP"} 1`,
				`http_requests_inflight{handler="test1",service="svc1",user_id="userVIP"} 1`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := prometheus.NewRegistry()
			test.config.Registry = reg
			mrecorder := libprometheus.NewRecorder(test.config)
			test.recordMetrics(mrecorder)

			// Get the metrics handler and serve.
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/metrics", nil)
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(rec, req)

			resp := rec.Result()

			// Check all metrics are present.
			if assert.Equal(t, http.StatusOK, resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)
				for _, expMetric := range test.expMetrics {
					assert.Contains(t, string(body), expMetric, "metric not present on the result")
				}
			}
		})
	}
}
