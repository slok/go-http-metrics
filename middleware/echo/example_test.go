package echo_test

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	echoMiddleware "github.com/slok/go-http-metrics/middleware/echo"
)

// EchoMiddleware shows how you would create a default middleware factory and use it
// to create an Echo compatible middleware.
func Example_echoMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our echo instance.
	e := echo.New()

	// Add our handler and middleware
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	}
	e.GET("/", h, echoMiddleware.Handler("", mdlw))

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", e); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
