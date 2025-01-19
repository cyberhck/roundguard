package rebalancer

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type LoadBalancer interface {
	ResetWithNewItems([]*httputil.ReverseProxy)
}

type Rebalancer[T LoadBalancer] struct {
	logger  *logrus.Logger
	proxies []*httputil.ReverseProxy
	lb      T
}

func (r *Rebalancer[T]) StartRebalancer(interval time.Duration, livenessEndpoint string) {
	ticker := time.NewTicker(interval)
	for {
		healthy, unhealthy := r.groupProxiesByLiveness(livenessEndpoint)
		if len(unhealthy) != 0 {
			r.logger.WithField("count", len(unhealthy)).Warn("Found some unhealthy processes")
		}
		r.lb.ResetWithNewItems(healthy)
		<-ticker.C
	}
}

func (r *Rebalancer[T]) GetLoadbalancer() T {
	return r.lb
}

func (r *Rebalancer[T]) groupProxiesByLiveness(livenessEndpoint string) ([]*httputil.ReverseProxy, []*httputil.ReverseProxy) {
	return lo.FilterReject(r.proxies, func(p *httputil.ReverseProxy, _ int) bool {
		request, err := http.NewRequest("GET", livenessEndpoint, nil)
		if err != nil {
			return false
		}
		resp, err := p.Transport.RoundTrip(request)
		if err != nil {
			return false
		}

		return resp.StatusCode == http.StatusOK
	})
}

func New[T LoadBalancer](proxies []*httputil.ReverseProxy, lb T, logger *logrus.Logger) *Rebalancer[T] {
	lb.ResetWithNewItems(proxies)

	return &Rebalancer[T]{
		lb:      lb,
		logger:  logger,
		proxies: proxies,
	}
}
