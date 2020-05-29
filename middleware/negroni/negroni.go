// Package negroni is a helper package to get a negroni compatible
// handler/middleware from the standard net/http Middleware factory.
package negroni

import (
	"net/http"

	"github.com/urfave/negroni"

	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

// Measure returns a Negroni measuring middleware from a Middleware instance.
func Measure(handlerID string, m middleware.Middleware) negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		std.Measure(handlerID, m, next).ServeHTTP(rw, r)
	})
}
