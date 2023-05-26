module github.com/slok/go-http-metrics/examples

go 1.19

replace (
	github.com/slok/go-http-metrics => ../
	github.com/slok/go-http-metrics/metrics/opencensus => ../metrics/opencensus
	github.com/slok/go-http-metrics/metrics/prometheus => ../metrics/prometheus
	github.com/slok/go-http-metrics/middleware/echo => ../middleware/echo
	github.com/slok/go-http-metrics/middleware/fasthttp => ../middleware/fasthttp
	github.com/slok/go-http-metrics/middleware/gin => ../middleware/gin
	github.com/slok/go-http-metrics/middleware/goji => ../middleware/goji
	github.com/slok/go-http-metrics/middleware/gorestful => ../middleware/gorestful
	github.com/slok/go-http-metrics/middleware/httprouter => ../middleware/httprouter
	github.com/slok/go-http-metrics/middleware/iris => ../middleware/iris
	github.com/slok/go-http-metrics/middleware/negroni => ../middleware/negroni
)

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.2
	github.com/emicklei/go-restful/v3 v3.10.2
	github.com/fasthttp/router v1.4.19
	github.com/gin-gonic/gin v1.9.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/justinas/alice v1.2.0
	github.com/kataras/iris/v12 v12.2.0
	github.com/labstack/echo/v4 v4.10.2
	github.com/prometheus/client_golang v1.15.1
	github.com/slok/go-http-metrics v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/metrics/opencensus v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/metrics/prometheus v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/echo v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/fasthttp v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/gin v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/goji v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/gorestful v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/httprouter v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/iris v0.0.0-00010101000000-000000000000
	github.com/slok/go-http-metrics/middleware/negroni v0.0.0-00010101000000-000000000000
	github.com/urfave/negroni v1.0.0
	github.com/valyala/fasthttp v1.47.0
	go.opencensus.io v0.24.0
	goji.io v2.0.2+incompatible
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet/v6 v6.2.0 // indirect
	github.com/Joker/jade v1.1.3 // indirect
	github.com/Shopify/goreferrer v0.0.0-20220729165902-8cddb4f5de06 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.9.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/flosch/pongo2/v4 v4.0.2 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/iris-contrib/schema v0.0.6 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kataras/blocks v0.0.7 // indirect
	github.com/kataras/golog v0.1.8 // indirect
	github.com/kataras/pio v0.0.11 // indirect
	github.com/kataras/sitemap v0.0.6 // indirect
	github.com/kataras/tunnel v0.0.4 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailgun/raymond/v2 v2.0.48 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.24 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.0 // indirect
	github.com/prometheus/statsd_exporter v0.23.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/schollz/closestmatch v2.1.0+incompatible // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	github.com/tdewolff/minify/v2 v2.12.5 // indirect
	github.com/tdewolff/parse/v2 v2.6.6 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/yosssi/ace v0.0.5 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
