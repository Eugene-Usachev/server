package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

var reuestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace: "GoServer",
	Subsystem: "http",
	Name:      "request",
	Objectives: map[float64]float64{
		0.5:  0.05,
		0.9:  0.01,
		0.99: 0.001,
	},
}, []string{"method", "path", "status"})

func ObserveRequest(duration float64, method string, path string, status int) {
	reuestMetrics.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration)
}
