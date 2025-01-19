package measurements

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var timer = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "request_duration_milliseconds",
	Buckets: []float64{10, 50, 200, 500, 800, 1_200, 2_000, 3_500, 5_000, 10_000},
}, []string{"scope", "method", "target", "status"})

func MeasureDuration(scope string, method string, target string) func(status string) time.Duration {
	startTime := time.Now()

	return func(status string) time.Duration {
		duration := time.Since(startTime)
		timer.WithLabelValues(
			scope,
			method,
			target,
			status,
		).Observe(float64(duration.Milliseconds()))

		return duration
	}
}
