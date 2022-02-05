package fasthttp_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promMetrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	fasthttpMiddleware "github.com/slok/go-http-metrics/middleware/fasthttp"
	"github.com/valyala/fasthttp"
)

func handleHello(rCtx *fasthttp.RequestCtx) {
	userID, ok := rCtx.UserValue("user_id").(string)
	if !ok {
		userID = "unknown"
	}

	rCtx.SetStatusCode(fasthttp.StatusOK)
	rCtx.SetBodyString(fmt.Sprintf("Hello, %s!", userID))
}

// FasthttpMiddleware shows how you would create a default middleware
// factory and use it to create a fasthttp compatible middleware.
func Example_fasthttpMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: promMetrics.NewRecorder(promMetrics.Config{}),
	})

	// Create our fasthttp instance.
	srv := &fasthttp.Server{
		Handler: fasthttpMiddleware.Handler("", mdlw, handleHello),
	}

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := srv.ListenAndServe(":8080"); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}

func Example_fasthttpCustomLabels() {
	mdlw := middleware.New(middleware.Config{
		Recorder: promMetrics.NewRecorder(promMetrics.Config{
			CustomLabels: []string{"user_id"},
		}),
	})

	mux := router.New()
	mux.GET("/{user_id}",
		func(c *fasthttp.RequestCtx) {
			mdlw.Measure("/hello", userIDReporter{c}, func() {
				handleHello(c)
			})
		},
	)

	srv := &fasthttp.Server{
		Handler: mux.Handler,
	}

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := srv.ListenAndServe(":8080"); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}

type userIDReporter struct {
	c *fasthttp.RequestCtx
}

func (r userIDReporter) Method() string {
	return string(r.c.Method())
}

func (r userIDReporter) Context() context.Context {
	return r.c
}

func (r userIDReporter) URLPath() string {
	return string(r.c.Path())
}

func (r userIDReporter) StatusCode() int {
	return r.c.Response.StatusCode()
}

func (r userIDReporter) BytesWritten() int64 {
	return int64(len(r.c.Response.Body()))
}

func (r userIDReporter) CustomLabels() []string {
	userID, ok := r.c.UserValue("user_id").(string)
	if !ok {
		return nil
	}

	return []string{userID}
}
