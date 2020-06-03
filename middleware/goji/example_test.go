package goji_test

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"goji.io"
	"goji.io/pat"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	gojimiddleware "github.com/slok/go-http-metrics/middleware/goji"
)

// GojiMiddleware shows how you would create a default middleware factory and use it
// to create a Goji compatible middleware.
func Example_gojiMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our Goji instance.
	mux := goji.NewMux()

	// Add our middleware.
	mux.Use(gojimiddleware.Handler("", mdlw))

	// Add our handler.
	mux.HandleFunc(pat.Get("/"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello world"))
	}))

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
