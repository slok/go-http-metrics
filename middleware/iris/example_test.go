package iris_test

import (
	"log"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	irismiddleware "github.com/slok/go-http-metrics/middleware/iris"
)

// IrisMiddleware shows how you would create a default middleware factory and use it
// to create an Iris compatible middleware.
func Example_irisMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our Iris Application instance.
	app := iris.New()

	// Add our handler and middleware
	h := func(ctx iris.Context) {
		ctx.WriteString("Hello world")
	}
	app.Get("/", irismiddleware.Handler("", mdlw), h)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := app.Listen(":8080", iris.WithOptimizations); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
