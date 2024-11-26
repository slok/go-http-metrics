// Package gorestful is a helper package to get a gorestful compatible middleware.
package gorestful

import (
	"context"

	gorestful "github.com/emicklei/go-restful/v3"

	"github.com/slok/go-http-metrics/middleware"
)

// Config is the configuration for the gorestful middleware.
type Config struct {
	UseRoutePath bool // When true, aggregate requests by their route path.  For example, "/users/{id}" instead of "/users/1", "/users/2", etc.
}

// Handler returns a gorestful measuring middleware with the default config.
func Handler(handlerID string, m middleware.Middleware) gorestful.FilterFunction {
	return HandlerWithConfig(handlerID, m, Config{})
}

// HandlerWithConfig returns a gorestful measuring middleware.
func HandlerWithConfig(handlerID string, m middleware.Middleware, config Config) gorestful.FilterFunction {
	return func(req *gorestful.Request, resp *gorestful.Response, chain *gorestful.FilterChain) {
		r := &reporter{req: req, resp: resp, config: config}
		m.Measure(handlerID, r, func() {
			chain.ProcessFilter(req, resp)
		})
	}
}

type reporter struct {
	req    *gorestful.Request
	resp   *gorestful.Response
	config Config
}

func (r *reporter) Method() string { return r.req.Request.Method }

func (r *reporter) Context() context.Context { return r.req.Request.Context() }

func (r *reporter) URLPath() string {
	if r.config.UseRoutePath {
		return r.req.SelectedRoutePath()
	}
	return r.req.Request.URL.Path
}

func (r *reporter) StatusCode() int { return r.resp.StatusCode() }

func (r *reporter) BytesWritten() int64 { return int64(r.resp.ContentLength()) }
