package web

import (
	"net/http"
	"time"

	"acortadorUrlService/components/logger"
	"acortadorUrlService/components/metrics"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := middleware.GetReqID(r.Context())
		if r.URL.Path == "/ping" {
			logger.LogDebug("Request started", "method", r.Method, "path", r.URL.Path, "request_id", reqID)
		} else {
			logger.LogInfo("Request started", "method", r.Method, "path", r.URL.Path, "request_id", reqID)
		}

		next.ServeHTTP(w, r)

		if r.URL.Path == "/ping" {
			logger.LogDebug("Request finished", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds(), "request_id", reqID)
		} else {
			logger.LogInfo("Request finished", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds(), "request_id", reqID)
		}
	})
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Milliseconds()

		var metricName string
		switch {
		case r.Method == "POST" && r.URL.Path == "/":
			metricName = metrics.MetricPostShortUrlDuration
		case r.Method == "DELETE" && r.URL.Path == "/":
			metricName = metrics.MetricDeleteShortUrlDuration
		case r.Method == "GET" && r.URL.Path != "/":
			metricName = metrics.MetricResolveShortUrlDuration
		default:
			return
		}
		metrics.PutDurationMetric(metricName, float64(duration))
	})
}
