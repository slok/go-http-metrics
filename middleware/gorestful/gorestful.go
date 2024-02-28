// Package gorestful is a helper package to get a gorestful compatible middleware.
package gorestful

import (
	"context"

	gorestful "github.com/emicklei/go-restful/v3"

	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a gorestful measuring middleware.
func Handler(handlerID string, m middleware.Middleware) gorestful.FilterFunction {
	return func(req *gorestful.Request, resp *gorestful.Response, chain *gorestful.FilterChain) {
		r := &reporter{req: req, resp: resp}
		m.Measure(handlerID, r, func() {
			chain.ProcessFilter(req, resp)
		})
	}
}

type reporter struct {
	req  *gorestful.Request
	resp *gorestful.Response
}

func (r *reporter) Method() string { return r.req.Request.Method }

func (r *reporter) Context() context.Context { return r.req.Request.Context() }

func (r *reporter) URLPath() string { return r.req.Request.URL.Path }

func (r *reporter) StatusCode() int { return r.resp.StatusCode() }

func (r *reporter) BytesWritten() int64 { return int64(r.resp.ContentLength()) }

func (r *reporter) CustomHeaders() map[string]string { return make(map[string]string) }
