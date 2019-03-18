package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
)

// This example shows how you could custom the middleware.
// It will use a prometheus recorder with a custom registry.
// It will set different configuration parameters to the middleware
// like a prefix or custom buckets for the histograms.
// It will set predefined handler ID to the handler middlewares so we maintain
// cardinality low instead of letting the middleware set the url path.
// If also groups the status codes.
func main() {
	// Create a custom registry for prometheus.
	reg := prometheus.NewRegistry()

	// Create our middleware.
	recorder := metrics.New(metrics.Config{
		Registry:        reg,
		Prefix:          "exampleapp",
		DurationBuckets: []float64{1, 2.5, 5, 10, 20, 40, 80, 160, 320, 640},
	})
	mdlw := middleware.New(middleware.Config{
		Recorder:      recorder,
		GroupedStatus: true,
	})

	// Create our server handlers.
	rooth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	testh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusAccepted) })
	testh2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	othetesth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusFound) })

	mux := http.NewServeMux()
	// Wrape our middleware on each of the different handlers with the ID of the handler
	// this way we reduce the cardinality, for example: `/test/2` and `/test/4` will
	// have the same `handler` label on the metric: `/test/:testID`
	mux.Handle("/", mdlw.Handler("/", rooth))
	mux.Handle("/test/1", mdlw.Handler("/test/:testID", testh))
	mux.Handle("/test/2", mdlw.Handler("/test/:testID", testh2))
	mux.Handle("/other-test", mdlw.Handler("/other-test", othetesth))

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, mux); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}()

	// Serve our metrics.
	go func() {
		log.Printf("metrics listening at %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.HandlerFor(reg, promhttp.HandlerOpts{})); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
