package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	fasthttprouter "github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	fasthttpmiddleware "github.com/slok/go-http-metrics/middleware/fasthttp"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
)

func main() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create fast HTTP router and global middleware.
	r := fasthttprouter.New()

	// Add our handler.
	r.GET("/", fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody([]byte("Hello world"))
	}))

	// Set up middleware.
	fasthttpHandler := fasthttpmiddleware.Handler("", mdlw, r.Handler)

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := fasthttp.ListenAndServe(srvAddr, fasthttpHandler); err != nil {
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
