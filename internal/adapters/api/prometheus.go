package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = prometheus.DefBuckets
)

const (
	patternReqsName     = "requests_total"
	patternLatencyName  = "requests_duration_seconds"
	httpServerSubsystem = "http_server"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewPatternMiddleware returns a new prometheus Middleware handler that groups requests by the chi routing pattern.
func NewPatternMiddleware(name string, buckets ...float64) func(next http.Handler) http.Handler {
	var m Middleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem:   httpServerSubsystem,
			Name:        patternReqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem:   httpServerSubsystem,
		Name:        patternLatencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path (with patterns).",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return m.patternHandler
}

func (c Middleware) patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		// It is pattern routing
		if rctx != nil {
			routePattern := strings.Join(rctx.RoutePatterns, "")
			routePattern = strings.ReplaceAll(routePattern, "/*/", "/")

			ww.Status()

			c.reqs.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, routePattern).Inc()
			c.latency.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, routePattern).Observe(time.Since(start).Seconds())
		} else {
			c.reqs.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, r.URL.Path).Inc()
			c.latency.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		}

	}
	return http.HandlerFunc(fn)
}
