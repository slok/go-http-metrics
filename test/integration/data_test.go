package integration

import "time"

var (
	expReqs = []handlerConfig{
		{Path: "/test/1", Method: "GET", Code: 201, ReturnData: "1", SleepDuration: 45 * time.Millisecond, NumberRequests: 10},
		{Path: "/test/2", Method: "POST", Code: 202, ReturnData: "22", SleepDuration: 95 * time.Millisecond, NumberRequests: 9},
		{Path: "/test/3", Method: "PATCH", Code: 203, ReturnData: "333", SleepDuration: 145 * time.Millisecond, NumberRequests: 8},
		{Path: "/test/4", Method: "DELETE", Code: 205, ReturnData: "4444", SleepDuration: 195 * time.Millisecond, NumberRequests: 7},
	}

	expMetrics = []string{
		`# HELP http_request_duration_seconds The latency of the HTTP requests.`,
		`# TYPE http_request_duration_seconds histogram`,
		`http_request_duration_seconds_bucket{code="201",handler="/test/1",method="GET",service="integration",le="0.05"} 10`,
		`http_request_duration_seconds_bucket{code="201",handler="/test/1",method="GET",service="integration",le="0.1"} 10`,
		`http_request_duration_seconds_bucket{code="201",handler="/test/1",method="GET",service="integration",le="0.15"} 10`,
		`http_request_duration_seconds_bucket{code="201",handler="/test/1",method="GET",service="integration",le="0.2"} 10`,
		`http_request_duration_seconds_bucket{code="201",handler="/test/1",method="GET",service="integration",le="+Inf"} 10`,
		`http_request_duration_seconds_count{code="201",handler="/test/1",method="GET",service="integration"} 10`,

		`http_request_duration_seconds_bucket{code="202",handler="/test/2",method="POST",service="integration",le="0.05"} 0`,
		`http_request_duration_seconds_bucket{code="202",handler="/test/2",method="POST",service="integration",le="0.1"} 9`,
		`http_request_duration_seconds_bucket{code="202",handler="/test/2",method="POST",service="integration",le="0.15"} 9`,
		`http_request_duration_seconds_bucket{code="202",handler="/test/2",method="POST",service="integration",le="0.2"} 9`,
		`http_request_duration_seconds_bucket{code="202",handler="/test/2",method="POST",service="integration",le="+Inf"} 9`,
		`http_request_duration_seconds_count{code="202",handler="/test/2",method="POST",service="integration"} 9`,

		`http_request_duration_seconds_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="0.05"} 0`,
		`http_request_duration_seconds_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="0.1"} 0`,
		`http_request_duration_seconds_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="0.15"} 8`,
		`http_request_duration_seconds_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="0.2"} 8`,
		`http_request_duration_seconds_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="+Inf"} 8`,
		`http_request_duration_seconds_count{code="203",handler="/test/3",method="PATCH",service="integration"} 8`,

		`http_request_duration_seconds_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="0.05"} 0`,
		`http_request_duration_seconds_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="0.1"} 0`,
		`http_request_duration_seconds_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="0.15"} 0`,
		`http_request_duration_seconds_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="0.2"} 7`,
		`http_request_duration_seconds_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="+Inf"} 7`,
		`http_request_duration_seconds_count{code="205",handler="/test/4",method="DELETE",service="integration"} 7`,

		`# HELP http_requests_inflight The number of inflight requests being handled at the same time.`,
		`# TYPE http_requests_inflight gauge`,
		`http_requests_inflight{handler="/test/1",service="integration"} 0`,
		`http_requests_inflight{handler="/test/2",service="integration"} 0`,
		`http_requests_inflight{handler="/test/3",service="integration"} 0`,
		`http_requests_inflight{handler="/test/4",service="integration"} 0`,

		`# HELP http_response_size_bytes The size of the HTTP responses.`,
		`# TYPE http_response_size_bytes histogram`,
		`http_response_size_bytes_bucket{code="201",handler="/test/1",method="GET",service="integration",le="1"} 10`,
		`http_response_size_bytes_bucket{code="201",handler="/test/1",method="GET",service="integration",le="2"} 10`,
		`http_response_size_bytes_bucket{code="201",handler="/test/1",method="GET",service="integration",le="3"} 10`,
		`http_response_size_bytes_bucket{code="201",handler="/test/1",method="GET",service="integration",le="4"} 10`,
		`http_response_size_bytes_bucket{code="201",handler="/test/1",method="GET",service="integration",le="5"} 10`,
		`http_response_size_bytes_bucket{code="201",handler="/test/1",method="GET",service="integration",le="+Inf"} 10`,
		`http_response_size_bytes_sum{code="201",handler="/test/1",method="GET",service="integration"} 10`,
		`http_response_size_bytes_count{code="201",handler="/test/1",method="GET",service="integration"} 10`,

		`http_response_size_bytes_bucket{code="202",handler="/test/2",method="POST",service="integration",le="1"} 0`,
		`http_response_size_bytes_bucket{code="202",handler="/test/2",method="POST",service="integration",le="2"} 9`,
		`http_response_size_bytes_bucket{code="202",handler="/test/2",method="POST",service="integration",le="3"} 9`,
		`http_response_size_bytes_bucket{code="202",handler="/test/2",method="POST",service="integration",le="4"} 9`,
		`http_response_size_bytes_bucket{code="202",handler="/test/2",method="POST",service="integration",le="5"} 9`,
		`http_response_size_bytes_bucket{code="202",handler="/test/2",method="POST",service="integration",le="+Inf"} 9`,
		`http_response_size_bytes_sum{code="202",handler="/test/2",method="POST",service="integration"} 18`,
		`http_response_size_bytes_count{code="202",handler="/test/2",method="POST",service="integration"} 9`,

		`http_response_size_bytes_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="1"} 0`,
		`http_response_size_bytes_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="2"} 0`,
		`http_response_size_bytes_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="3"} 8`,
		`http_response_size_bytes_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="4"} 8`,
		`http_response_size_bytes_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="5"} 8`,
		`http_response_size_bytes_bucket{code="203",handler="/test/3",method="PATCH",service="integration",le="+Inf"} 8`,
		`http_response_size_bytes_sum{code="203",handler="/test/3",method="PATCH",service="integration"} 24`,
		`http_response_size_bytes_count{code="203",handler="/test/3",method="PATCH",service="integration"} 8`,

		`http_response_size_bytes_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="1"} 0`,
		`http_response_size_bytes_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="2"} 0`,
		`http_response_size_bytes_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="3"} 0`,
		`http_response_size_bytes_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="4"} 7`,
		`http_response_size_bytes_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="5"} 7`,
		`http_response_size_bytes_bucket{code="205",handler="/test/4",method="DELETE",service="integration",le="+Inf"} 7`,
		`http_response_size_bytes_sum{code="205",handler="/test/4",method="DELETE",service="integration"} 28`,
		`http_response_size_bytes_count{code="205",handler="/test/4",method="DELETE",service="integration"} 7`,
	}
)
