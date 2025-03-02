package monitoring

import "github.com/prometheus/client_golang/prometheus"

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests processed",
	},
	[]string{"method", "endpoint"},
)

func InitPrometheus() {
	prometheus.MustRegister(requestCount)
}

func GetRequestCount() *prometheus.CounterVec {
	return requestCount
}
