package gorestful_test

import (
	"log"
	"net/http"

	gorestful "github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	gorestfulmiddleware "github.com/slok/go-http-metrics/middleware/gorestful"
)

// GorestfulMiddleware shows how you would create a default middleware factory and use it
// to create a Gorestful compatible middleware.
func Example_gorestfulMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our gorestful instance.
	c := gorestful.NewContainer()

	// Add the middleware for all routes.
	c.Filter(gorestfulmiddleware.Handler("", mdlw))

	// Add our handler,
	ws := &gorestful.WebService{}
	ws.Route(ws.GET("/").To(func(_ *gorestful.Request, resp *gorestful.Response) {
		_ = resp.WriteEntity("Hello world")
	}))
	c.Add(ws)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", c); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
