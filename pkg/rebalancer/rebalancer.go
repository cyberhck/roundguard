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

type Rebalancer struct {
	logger  *logrus.Logger
	proxies []*httputil.ReverseProxy
	lb      LoadBalancer
}

func (r *Rebalancer) StartRebalancer() {
	for {
		time.Sleep(time.Second)
		healthy, unhealthy := lo.FilterReject(r.proxies, func(p *httputil.ReverseProxy, _ int) bool {
			request, err := http.NewRequest("GET", "/live", nil)
			if err != nil {
				return false
			}
			resp, err := p.Transport.RoundTrip(request)
			if err != nil {
				return false
			}

			return resp.StatusCode == http.StatusOK
		})
		if len(unhealthy) != 0 {
			r.logger.WithField("count", len(unhealthy)).Warn("Found some unhealthy processes")
		}
		r.lb.ResetWithNewItems(healthy)
	}
}

func New(proxies []*httputil.ReverseProxy, lb LoadBalancer, logger *logrus.Logger) *Rebalancer {
	lb.ResetWithNewItems(proxies)

	return &Rebalancer{
		lb:      lb,
		logger:  logger,
		proxies: proxies,
	}
}
