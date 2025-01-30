// Package fasthttp is a helper package to get a fasthttp compatible middleware.
package fasthttp

import (
	"context"

	"github.com/slok/go-http-metrics/middleware"
	"github.com/valyala/fasthttp"
)

// Handler returns a fasthttp measuring middleware.
func Handler(handlerID string, m middleware.Middleware, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		m.Measure(handlerID, reporter{c}, func() {
			next(c)
		})
	}
}

type reporter struct {
	c *fasthttp.RequestCtx
}

func (r reporter) Method() string {
	return string(r.c.Method())
}

func (r reporter) Context() context.Context {
	return r.c
}

func (r reporter) URLPath() string {
	return string(r.c.Path())
}

func (r reporter) StatusCode() int {
	return r.c.Response.StatusCode()
}

func (r reporter) BytesWritten() int64 {
	return int64(len(r.c.Response.Body()))
}

func (r reporter) CustomHeaders() map[string]string { return make(map[string]string) }
