// Package gin is a helper package to get a gin compatible
// handler/middleware from the standard net/http Middleware factory.
package gin

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/slok/go-http-metrics/middleware"
)

// Measure returns a Gin measure middleware.
func Measure(handlerID string, m middleware.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := &reporter{c: c}
		m.Measure(handlerID, r, func() {
			c.Next()
		})
	}
}

type reporter struct {
	c *gin.Context
}

func (r *reporter) Method() string { return r.c.Request.Method }

func (r *reporter) Context() context.Context { return r.c.Request.Context() }

func (r *reporter) URLPath() string { return r.c.Request.URL.Path }

func (r *reporter) StatusCode() int { return r.c.Writer.Status() }

func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }
