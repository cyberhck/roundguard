package loadbalancer

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/cyberhck/roundguard/pkg/http/types"
	"github.com/cyberhck/roundguard/pkg/rebalancer"
	"github.com/cyberhck/roundguard/pkg/roundrobin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func CreateServer(addr string, rb *rebalancer.Rebalancer[*roundrobin.RoundRobin[*httputil.ReverseProxy]], logger *logrus.Logger) (*http.Server, error) {
	router := http.NewServeMux()
	router.Handle("GET /metrics", promhttp.Handler())
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy, err := rb.GetLoadbalancer().GetNext()
		if err != nil {
			logger.WithField("error", err).Error("No proxy available")
			w.WriteHeader(http.StatusServiceUnavailable)
			err = json.NewEncoder(w).Encode(types.ErrorResponse{
				Title:       "No healthy upstream",
				Description: "There's no instance of application that can process this request",
			})
			if err != nil {
				logger.WithField("error", err).Error("Failed to write response")
			}
			return
		}
		(*proxy).ServeHTTP(w, r)
	})

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}, nil
}
