package users

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus Metrics
var (
	totalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_api_requests_total",
		Help: "The current total requests made to the auth api",
	}, []string{"code", "method", "route"})
	buckets       = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	responseTimes = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "namespace",
		Name:      "auth_api_request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
		Buckets:   buckets,
	}, []string{"route", "method", "status_code"})
)

func registerMetrics() error {
	err := prometheus.Register(totalRequests)
	if err != nil {
		log.Printf("failed to register total requests metric err=\"%s\"", err)
		return err
	}

	err = prometheus.Register(responseTimes)
	if err != nil {
		log.Printf("failed to register response times metric err=\"%s\"", err)
		return err
	}
	return nil
}
