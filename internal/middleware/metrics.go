package middleware

import (
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "route", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Histogram of response time for handler",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "route"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        recorder := &statusRecorder{ResponseWriter: w, status: 200}
        next.ServeHTTP(recorder, r)

        routePattern := chi.RouteContext(r.Context()).RoutePattern()

        httpRequestsTotal.WithLabelValues(r.Method, routePattern, http.StatusText(recorder.status)).Inc()
        httpRequestDuration.WithLabelValues(r.Method, routePattern).Observe(time.Since(start).Seconds())
    })
}
