// Package gin is a helper package to get a gin compatible
// handler/middleware from the standard net/http Middleware factory.
package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Gin compatible middleware from a Middleware factory instance.
// The first handlerID argument is the same argument passed on Middleware.Handler method.
func Handler(handlerID string, m middleware.Middleware) gin.HandlerFunc {
	// Create a dummy handler to wrap the middleware chain of Gin, this way Middleware
	// interface can wrap the Gin chain.
	return func(c *gin.Context) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Writer = &ginResponseWriter{
				ResponseWriter: c.Writer,
				middlewareRW:   w,
			}
			c.Next()
		})
		m.Handler(handlerID, h).ServeHTTP(c.Writer, c.Request)
	}
}

// ginResponseWriter is a helper type that intercepts the middleware ResponseWriter
// interceptor.
// This is required because gin's context Writer (c.Writer) is a gin.ResponseWriter
// interface and we can't access to the internal object http.ResponseWriter, so
// we already know that our middleware intercepts the regular http.ResponseWriter,
// and doesn't change anything, just intercepts to read information. So in order to
// get this information on our interceptor we create a gin.ResponseWriter implementation
// that will call the real gin.Context.Writer and our interceptor. This way Gin gets the
// information and our interceptor also.
type ginResponseWriter struct {
	middlewareRW http.ResponseWriter
	gin.ResponseWriter
}

func (w *ginResponseWriter) WriteHeader(statusCode int) {
	w.middlewareRW.WriteHeader(statusCode)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ginResponseWriter) Write(p []byte) (int, error) {
	w.middlewareRW.Write(p)
	return w.ResponseWriter.Write(p)
}
