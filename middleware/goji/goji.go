// Package goji is a helper package to get a goji compatible middleware.
package goji

import (
	"net/http"

	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

// Handler returns a Goji measuring middleware.
func Handler(handlerID string, m middleware.Middleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return std.Handler(handlerID, m, next)
	}
}
