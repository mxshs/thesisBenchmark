package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var TTS prometheus.Histogram
var TotalRequests prometheus.Counter

func SetMetrics() {
	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_req_2_tickers",
			Help: "Total client requests to tickers RPC",
		})

	TTS = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "time_2_serve_ticker",
			Help:    "Time needed to serve a request to tickers RPC",
			Buckets: prometheus.LinearBuckets(0, 1, 20),
		})
}
