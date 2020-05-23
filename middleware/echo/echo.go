// Package echo is a helper package to get an echo compatible
// handler/middleware from the standard net/http Middleware factory.
package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Echo compatible middleware from a Middleware factory instance.
// The first handlerID argument is the same argument passed on Middleware.Handler method.
func Handler(handlerID string, m middleware.Middleware) echo.MiddlewareFunc {
	// Wrap wrapping handler with echo's WrapMiddleware helper
	return echo.WrapMiddleware(func(next http.Handler) http.Handler {
		return m.Handler(handlerID, next)
	})
}
