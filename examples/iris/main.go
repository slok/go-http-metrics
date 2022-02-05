package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	irismiddleware "github.com/slok/go-http-metrics/middleware/iris"
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

	// Create Iris engine and global middleware.
	app := iris.New()
	app.Use(irismiddleware.Handler("", mdlw))

	// Add our handler.
	app.Get("/", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusAccepted)
		_, _ = ctx.WriteString("Hello world!")
	})

	app.Get("/json", func(ctx iris.Context) {
		ctx.JSON(map[string]string{"hello": "world"}) // nolint: errcheck
	})

	app.Get("/wrong", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusTooManyRequests)
		_, _ = ctx.WriteString("oops")

	})

	err := app.Build()
	if err != nil {
		panic(err)
	}

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, app); err != nil {
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
