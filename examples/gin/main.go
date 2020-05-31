package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
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

	// Create Gin engine and global middleware.
	engine := gin.New()
	engine.Use(ginmiddleware.Handler("", mdlw))

	// Add our handler.
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello %s", "world")
	})
	engine.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, map[string]string{"hello": "world"})
	})
	engine.GET("/yaml", func(c *gin.Context) {
		c.YAML(http.StatusCreated, map[string]string{"hello": "world"})
	})
	engine.GET("/wrong", func(c *gin.Context) {
		c.String(http.StatusTooManyRequests, "oops")
	})

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, engine); err != nil {
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
