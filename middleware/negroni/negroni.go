// Package negroni is a helper package to get a negroni compatible
// handler/middleware from the standard net/http Middleware factory.
package negroni

import (
	"net/http"

	"github.com/urfave/negroni"

	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Negroni compatible middleware from a Middleware factory instance.
// The first handlerID argument is the same argument passed on Middleware.Handler method.
func Handler(handlerID string, m middleware.Middleware) negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		m.Handler(handlerID, next).ServeHTTP(rw, r)
	})
}
