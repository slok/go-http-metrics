// Package httprouter is a helper package to get a httprouter compatible
// handler/middleware from the standatd net/http Middleware factory.
package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

// Measure returns a httprouter.Handler measured middlware.
func Measure(handlerID string, next httprouter.Handle, m middleware.Middleware) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Dummy handler to wrap httprouter Handle type
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, p)
		})

		std.Measure(handlerID, m, h).ServeHTTP(w, r)
	}
}
