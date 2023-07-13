// Package iris is a helper package to get an Iris compatible middleware.
package iris

import (
	"context"

	"github.com/kataras/iris/v12"

	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Iris measuring middleware.
func Handler(handlerID string, m middleware.Middleware) iris.Handler {
	return func(ctx iris.Context) {
		r := &reporter{ctx: ctx}
		m.Measure(handlerID, r, ctx, func() {
			ctx.Next()
		})
	}
}

type reporter struct {
	ctx iris.Context
}

func (r *reporter) Method() string { return r.ctx.Method() }

func (r *reporter) Context() context.Context { return r.ctx.Request().Context() }

func (r *reporter) URLPath() string { return r.ctx.Path() }

func (r *reporter) StatusCode() int { return r.ctx.GetStatusCode() }

func (r *reporter) BytesWritten() int64 { return int64(r.ctx.ResponseWriter().Written()) }
