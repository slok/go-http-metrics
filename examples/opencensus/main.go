package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ocprometheus "contrib.go.opencensus.io/exporter/prometheus"
	ocmmetrics "github.com/slok/go-http-metrics/metrics/opencensus"
	"go.opencensus.io/stats/view"

	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
)

// This example will show how you could use opencensus with go-http-middleware
// and serve the metrics in Prometheus format (through OpenCensus).
func main() {
	// Create OpenCensus with Prometheus.
	ocexporter, err := ocprometheus.NewExporter(ocprometheus.Options{})
	if err != nil {
		log.Panicf("error creating OpenCensus exporter: %s", err)
	}
	view.RegisterExporter(ocexporter)
	rec, err := ocmmetrics.NewRecorder(ocmmetrics.Config{})
	if err != nil {
		log.Panicf("error creating OpenCensus metrics recorder: %s", err)
	}

	// Create our middleware.
	mdlw := middleware.New(middleware.Config{
		Recorder: rec,
	})

	// Create our server.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })
	mux.HandleFunc("/test1/test2", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusAccepted) })
	mux.HandleFunc("/test1/test4", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNonAuthoritativeInfo) })
	mux.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("/test3", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusResetContent) })

	// Wrap our main handler, we pass empty handler ID so the middleware inferes
	// the handler label from the URL.
	h := std.Handler("", mdlw, mux)

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, h); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}()

	// Serve our metrics.
	go func() {
		log.Printf("metrics listening at %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, ocexporter); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
