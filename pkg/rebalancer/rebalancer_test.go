//go:generate go run go.uber.org/mock/mockgen -destination=../../mocks/mock_round_robin.go -package mocks -typed github.com/cyberhck/roundguard/pkg/rebalancer LoadBalancer
package rebalancer_test

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/cyberhck/roundguard/mocks"
	"github.com/cyberhck/roundguard/pkg/rebalancer"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	t.Run("it should remove unhealthy proxies in a given duration", func(t *testing.T) {
		mockLb := mocks.NewMockLoadBalancer(gomock.NewController(t))
		count := 0
		mockLb.EXPECT().ResetWithNewItems(gomock.Any()).DoAndReturn(func(proxies []*httputil.ReverseProxy) {
			if count == 0 {
				assert.Len(t, proxies, 2)
			}
			if count == 1 {
				assert.Len(t, proxies, 0)
			}
			count++
		}).AnyTimes()
		tr := httpmock.NewMockTransport()
		tr.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
		p1 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p1.Transport = tr
		p2 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p2.Transport = tr
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		rb := rebalancer.New([]*httputil.ReverseProxy{p1, p2}, mockLb, logger)
		go rb.StartRebalancer()
		time.Sleep(2 * time.Second)
	})
}
