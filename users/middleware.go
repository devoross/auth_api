package users

import (
	"net/http"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := statusRecorder{w, 200}

		next.ServeHTTP(&rec, r)

		duration := time.Since(startTime).Seconds()

		statusCode := strconv.Itoa(rec.statusCode)
		totalRequests.WithLabelValues(statusCode, r.Method, r.RequestURI).Inc()
		responseTimes.WithLabelValues(r.RequestURI, r.Method, statusCode).Observe(duration)
	})
}
