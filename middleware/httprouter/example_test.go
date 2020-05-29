package httprouter_test

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
	stdmiddleware "github.com/slok/go-http-metrics/middleware/std"
)

// HTTPRouterMiddlewareByHandler shows how you would create a default middleware factory
// and use it to create a httprouter compatible middleware setting by handler instead of
// main router.
func Example_httprouterMiddlewareByHandler() {
	// Create our handler.
	myHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world! " + id))
	}

	// Create our router.
	r := httprouter.New()

	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	// Add the middleware.
	r.GET("/test/:id", httproutermiddleware.Measure("/test/:id", myHandler, mdlw))
	r.GET("/test2/:id", httproutermiddleware.Measure("/test2/:id", myHandler, mdlw))

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}

// HTTPRouterMiddlewareOnRouter shows how you would create a default middleware factory
// and use it to wrapon httprouter Router (that satisfies http.Handler interface).
func Example_httprouterMiddlewareOnRouter() {
	// Create our handler.
	myHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world! " + id))
	}

	// Create our router and add the middleware.
	r := httprouter.New()

	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	r.GET("/test/:id", myHandler)
	r.GET("/test2/:id", myHandler)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Wrap the router with the middleware.
	h := stdmiddleware.Measure("", mdlw, r)

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
