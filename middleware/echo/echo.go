// Package echo is a helper package to get an echo compatible
// handler/middleware from the standard net/http Middleware factory.
package echo

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Echo measuring middleware.
func Handler(handlerID string, m middleware.Middleware) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			r := &reporter{c: c}
			var err error
			m.Measure(handlerID, r, func() {
				err = h(c)
			})
			return err
		})
	}
}

type reporter struct {
	c echo.Context
}

func (r *reporter) Method() string { return r.c.Request().Method }

func (r *reporter) Context() context.Context { return r.c.Request().Context() }

func (r *reporter) URLPath() string { return r.c.Request().URL.Path }

func (r *reporter) StatusCode() int { return r.c.Response().Status }

func (r *reporter) BytesWritten() int64 { return r.c.Response().Size }
