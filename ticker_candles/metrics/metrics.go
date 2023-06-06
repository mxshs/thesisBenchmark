package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var TTS prometheus.Histogram
var TotalRequests prometheus.Counter

func SetMetrics() {

	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_requests_processed",
			Help: "Total client requests to tickers RPC",
		})

	TTS = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "request_processing_time",
			Help:    "Time needed to serve a request to tickers RPC",
			Buckets: prometheus.LinearBuckets(0, 1, 20),
		})
}
