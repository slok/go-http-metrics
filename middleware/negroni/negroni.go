// Package negroni is a helper package to get a negroni compatible middleware.
package negroni

import (
	"net/http"

	"github.com/urfave/negroni"

	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

// Handler returns a Negroni measuring middleware.
func Handler(handlerID string, m middleware.Middleware) negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		std.Handler(handlerID, m, next).ServeHTTP(rw, r)
	})
}
