package mux_test

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	muxmiddleware "github.com/slok/go-http-metrics/middleware/mux"
)

// MuxMiddleware shows how you would create a default middleware factory and use it
// to create a Gorilla Mux `http.Handler` compatible middleware.
func Example_muxMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our handler.
	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world!"))
	})

	// Wrap our handler with the middleware.
	h := muxmiddleware.Handler("", mdlw, myHandler)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
