package integration

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gorestful "github.com/emicklei/go-restful/v3"
	fasthttprouter "github.com/fasthttp/router"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kataras/iris/v12"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
	"github.com/valyala/fasthttp"
	"goji.io"
	"goji.io/pat"

	metricsprometheus "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	echomiddleware "github.com/slok/go-http-metrics/middleware/echo"
	fasthttpmiddleware "github.com/slok/go-http-metrics/middleware/fasthttp"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	gojimiddleware "github.com/slok/go-http-metrics/middleware/goji"
	gorestfulmiddleware "github.com/slok/go-http-metrics/middleware/gorestful"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
	irismiddleware "github.com/slok/go-http-metrics/middleware/iris"
	negronimiddleware "github.com/slok/go-http-metrics/middleware/negroni"
	stdmiddleware "github.com/slok/go-http-metrics/middleware/std"
)

// server is the interface used by the integration tests to return a listening server
// where we can make request and then test the metrics.
type server interface {
	Close()
	URL() string
}

type testServer struct{ server *httptest.Server }

func (t testServer) Close()      { t.server.Close() }
func (t testServer) URL() string { return t.server.URL }

type netListenerServer struct {
	ln net.Listener
}

func (n netListenerServer) Close()      { _ = n.ln.Close() }
func (n netListenerServer) URL() string { return "http://" + n.ln.Addr().String() }

// handlerConfig is the configuration the servers will need to set up to properly
// execute the tests.
type handlerConfig struct {
	Path           string
	Code           int
	Method         string
	ReturnData     string
	SleepDuration  time.Duration
	NumberRequests int
}

func TestMiddlewarePrometheus(t *testing.T) {
	tests := map[string]struct {
		server func(m middleware.Middleware, hc []handlerConfig) server
	}{
		"STD http.Handler": {server: prepareHandlerSTD},
		"Negroni":          {server: prepareHandlerNegroni},
		"HTTPRouter":       {server: prepareHandlerHTTPRouter},
		"Gorestful":        {server: prepareHandlerGorestful},
		"Gin":              {server: prepareHandlerGin},
		"Echo":             {server: prepareHandlerEcho},
		"Goji":             {server: prepareHandlerGoji},
		"Chi":              {server: prepareHandlerChi},
		"Alice":            {server: prepareHandlerAlice},
		"Gorilla":          {server: prepareHandlerGorilla},
		"Fasthttp":         {server: prepareHandlerFastHTTP},
		"Iris":             {server: prepareHandlerIris},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			reg := prometheus.NewRegistry()
			rec := metricsprometheus.NewRecorder(metricsprometheus.Config{
				Registry:        reg,
				DurationBuckets: []float64{0.05, 0.1, 0.15, 0.2},
				SizeBuckets:     []float64{1, 2, 3, 4, 5},
			})
			mdlw := middleware.New(middleware.Config{
				Service:  "integration",
				Recorder: rec,
			})

			server := test.server(mdlw, expReqs)
			metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

			// Test.
			testMiddlewareRequests(t, server, expReqs)
			testMiddlewarePrometheusMetrics(t, metricsHandler, expMetrics)
		})
	}
}

func testMiddlewareRequests(t *testing.T, server server, expReqs []handlerConfig) {
	require := require.New(t)
	assert := assert.New(t)

	// Setup server cleanup.
	t.Cleanup(func() { server.Close() })

	// Make all the requests.
	for _, config := range expReqs {
		for i := 0; i < config.NumberRequests; i++ {
			r, err := http.NewRequest(config.Method, server.URL()+config.Path, nil)
			require.NoError(err)
			resp, err := http.DefaultClient.Do(r)
			require.NoError(err)

			// Check.
			assert.Equal(config.Code, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			require.NoError(err)
			assert.Equal(config.ReturnData, string(b))
		}
	}
}

func testMiddlewarePrometheusMetrics(t *testing.T, h http.Handler, expMetrics []string) {
	require := require.New(t)
	assert := assert.New(t)

	// Setup server.
	server := httptest.NewServer(h)
	t.Cleanup(func() { server.Close() })

	// Get metrics.
	r, err := http.NewRequest(http.MethodGet, server.URL+"/metrics", nil)
	require.NoError(err)
	resp, err := http.DefaultClient.Do(r)
	require.NoError(err)

	// Check.
	b, err := io.ReadAll(resp.Body)
	require.NoError(err)
	metrics := string(b)

	// Make all the requests.
	for _, expMetric := range expMetrics {
		assert.Contains(metrics, expMetric)
	}
}

func prepareHandlerSTD(m middleware.Middleware, hc []handlerConfig) server {
	// Setup handlers.
	mux := http.NewServeMux()
	for _, h := range hc {
		h := h
		mux.Handle(h.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != h.Method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			time.Sleep(h.SleepDuration)
			w.WriteHeader(h.Code)
			w.Write([]byte(h.ReturnData)) // nolint: errcheck
		}))
	}

	// Setup server and middleware.
	h := stdmiddleware.Handler("", m, mux)

	return testServer{server: httptest.NewServer(h)}
}

func prepareHandlerNegroni(m middleware.Middleware, hc []handlerConfig) server {
	// Setup handlers.
	mux := http.NewServeMux()
	for _, h := range hc {
		h := h
		mux.Handle(h.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != h.Method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			time.Sleep(h.SleepDuration)
			w.WriteHeader(h.Code)
			w.Write([]byte(h.ReturnData)) // nolint: errcheck
		}))
	}

	// Setup server and middleware.
	n := negroni.Classic()
	n.Use(negronimiddleware.Handler("", m))
	n.UseHandler(mux)

	return testServer{server: httptest.NewServer(n)}
}

func prepareHandlerHTTPRouter(m middleware.Middleware, hc []handlerConfig) server {
	r := httprouter.New()

	// Setup handlers.
	for _, h := range hc {
		h := h
		hr := func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			time.Sleep(h.SleepDuration)
			w.WriteHeader(h.Code)
			w.Write([]byte(h.ReturnData)) // nolint: errcheck
		}

		// Setup middleware on each of the routes.
		r.Handle(h.Method, h.Path, httproutermiddleware.Handler("", hr, m))
	}

	return testServer{server: httptest.NewServer(r)}
}

func prepareHandlerGorestful(m middleware.Middleware, hc []handlerConfig) server {
	// Setup server and middleware.
	c := gorestful.NewContainer()
	c.Filter(gorestfulmiddleware.Handler("", m))

	// Setup handlers.
	ws := &gorestful.WebService{}
	for _, h := range hc {
		h := h
		ws.Route(ws.Method(h.Method).Path(h.Path).To(func(_ *gorestful.Request, resp *gorestful.Response) {
			time.Sleep(h.SleepDuration)
			resp.WriteHeader(h.Code)
			resp.Write([]byte(h.ReturnData)) // nolint: errcheck
		}))
	}
	c.Add(ws)

	return testServer{server: httptest.NewServer(c)}
}

func prepareHandlerGin(m middleware.Middleware, hc []handlerConfig) server {
	// Setup server and middleware.
	e := gin.New()
	e.Use(ginmiddleware.Handler("", m))

	// Setup handlers.
	for _, h := range hc {
		h := h
		e.Handle(h.Method, h.Path, func(c *gin.Context) {
			time.Sleep(h.SleepDuration)
			c.String(h.Code, h.ReturnData)
		})
	}

	return testServer{server: httptest.NewServer(e)}
}

func prepareHandlerEcho(m middleware.Middleware, hc []handlerConfig) server {
	// Setup server and middleware.
	e := echo.New()
	e.Use(echomiddleware.Handler("", m))

	// Setup handlers.
	for _, h := range hc {
		h := h
		e.Add(h.Method, h.Path, func(c echo.Context) error {
			time.Sleep(h.SleepDuration)
			return c.String(h.Code, h.ReturnData)
		})
	}

	return testServer{server: httptest.NewServer(e)}
}

func prepareHandlerGoji(m middleware.Middleware, hc []handlerConfig) server {
	// Setup server and middleware.
	mux := goji.NewMux()
	mux.Use(gojimiddleware.Handler("", m))

	// Setup handlers.
	for _, h := range hc {
		h := h
		mux.HandleFunc(pat.NewWithMethods(h.Path, h.Method), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(h.SleepDuration)
			w.WriteHeader(h.Code)
			w.Write([]byte(h.ReturnData)) // nolint: errcheck
		}))
	}

	return testServer{server: httptest.NewServer(mux)}
}

func prepareHandlerChi(m middleware.Middleware, hc []handlerConfig) server {
	// Setup server and middleware.
	mux := chi.NewMux()
	mux.Use(stdmiddleware.HandlerProvider("", m))

	// Setup handlers.
	for _, h := range hc {
		h := h
		mux.Method(h.Method, h.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(h.SleepDuration)
			w.WriteHeader(h.Code)
			w.Write([]byte(h.ReturnData)) // nolint: errcheck
		}))
	}

	return testServer{server: httptest.NewServer(mux)}
}

func prepareHandlerAlice(m middleware.Middleware, hc []handlerConfig) server {
	// Setup handlers.
	mux := http.NewServeMux()
	for _, h := range hc {
		h := h
		mux.Handle(h.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != h.Method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			time.Sleep(h.SleepDuration)
			w.WriteHeader(h.Code)
			w.Write([]byte(h.ReturnData)) // nolint: errcheck
		}))
	}

	// Setup server and middleware.
	h := alice.New(stdmiddleware.HandlerProvider("", m)).Then(mux)

	return testServer{server: httptest.NewServer(h)}
}

func prepareHandlerGorilla(m middleware.Middleware, hc []handlerConfig) server {
	// Setup handlers.
	r := mux.NewRouter()
	for _, h := range hc {
		h := h
		r.Methods(h.Method).
			Path(h.Path).
			HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(h.SleepDuration)
				w.WriteHeader(h.Code)
				w.Write([]byte(h.ReturnData)) // nolint: errcheck
			}))
	}

	// Setup middleware.
	r.Use(stdmiddleware.HandlerProvider("", m))

	return testServer{server: httptest.NewServer(r)}
}

func prepareHandlerFastHTTP(m middleware.Middleware, hc []handlerConfig) server {
	// Setup handlers.
	r := fasthttprouter.New()
	for _, h := range hc {
		h := h
		r.Handle(h.Method, h.Path, fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
			time.Sleep(h.SleepDuration)
			ctx.SetStatusCode(h.Code)
			ctx.SetBody([]byte(h.ReturnData))
		}))
	}

	// Setup middleware.
	fasthttpHandler := fasthttpmiddleware.Handler("", m, r.Handler)

	// Setup server.
	// fasthttp doesn't use the regular http std lib server, so we need to use
	// a custom net TCP listener so we can obtain the random port URL.
	ln, _ := net.Listen("tcp", "127.0.0.1:0") // `:0` for random port.
	go func() {
		fasthttp.Serve(ln, fasthttpHandler) // nolint: errcheck
	}()

	return netListenerServer{ln: ln}
}

func prepareHandlerIris(m middleware.Middleware, hc []handlerConfig) server {
	// Setup server and middleware.
	app := iris.New()
	app.Use(irismiddleware.Handler("", m))

	// Set handlers.
	for _, h := range hc {
		h := h
		app.Handle(h.Method, h.Path, iris.Handler(func(ctx iris.Context) {
			time.Sleep(h.SleepDuration)
			ctx.StatusCode(h.Code)
			ctx.WriteString(h.ReturnData) // nolint: errcheck
		}))
	}

	app.Build() // nolint: errcheck

	return testServer{server: httptest.NewServer(app)}
}
