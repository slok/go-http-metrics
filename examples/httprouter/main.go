package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
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

	// Create our handlers.
	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("root"))
	}
	h1 := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test " + id))
	}

	h2 := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("test2 " + id))
	}

	// Create our router.
	r := httprouter.New()

	// Add the middleware to each route.
	r.GET("/", httproutermiddleware.Handler("/", h, mdlw))
	r.GET("/test/:id", httproutermiddleware.Handler("/test/:id", h1, mdlw))
	r.GET("/test2/:id", httproutermiddleware.Handler("/test2/:id", h2, mdlw))

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, r); err != nil {
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
