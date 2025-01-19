package measurements

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var availableProxies = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "proxy_count",
}, []string{"availability"})

func MeasureAvailableProxies(available int, unavailable int) {
	availableProxies.WithLabelValues("available").Set(float64(available))
	availableProxies.WithLabelValues("unavailable").Set(float64(unavailable))
}
