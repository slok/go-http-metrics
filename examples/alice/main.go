package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
)

func main() {
	// Create our middleware.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
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

	// Wrap with middleware.
	h := alice.New(std.HandlerProvider("", mdlw)).Then(mux)

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
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
