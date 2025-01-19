//go:generate go run go.uber.org/mock/mockgen -destination=../../mocks/mock_round_robin.go -package mocks -typed github.com/cyberhck/roundguard/pkg/rebalancer LoadBalancer
package rebalancer_test

import (
	"io"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/cyberhck/roundguard/pkg/roundrobin"
	"github.com/stretchr/testify/assert"

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
				assert.Len(t, proxies, 3)
			}
			if count == 1 {
				assert.Len(t, proxies, 1)
			}
			count++
		}).AnyTimes()
		tr1 := httpmock.NewMockTransport()
		tr1.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
		p1 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p1.Transport = tr1
		tr2 := httpmock.NewMockTransport()
		tr2.RegisterNoResponder(httpmock.NewStringResponder(200, ""))
		p2 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p2.Transport = tr2
		tr3 := httpmock.NewMockTransport()
		p3 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p3.Transport = tr3
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		rb := rebalancer.New([]*httputil.ReverseProxy{p1, p2, p3}, mockLb, logger)
		go rb.StartRebalancer(time.Millisecond*200, "/live")
		time.Sleep(time.Microsecond * 250)
	})
	t.Run("if liveness endpoint is invalid, it removes everything", func(t *testing.T) {
		mockLb := mocks.NewMockLoadBalancer(gomock.NewController(t))
		count := 0
		mockLb.EXPECT().ResetWithNewItems(gomock.Any()).DoAndReturn(func(proxies []*httputil.ReverseProxy) {
			if count == 0 {
				assert.Len(t, proxies, 3)
			}
			if count == 1 {
				assert.Len(t, proxies, 0)
			}
			count++
		}).AnyTimes()
		tr1 := httpmock.NewMockTransport()
		tr1.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
		p1 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p1.Transport = tr1
		tr2 := httpmock.NewMockTransport()
		tr2.RegisterNoResponder(httpmock.NewStringResponder(200, ""))
		p2 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p2.Transport = tr2
		tr3 := httpmock.NewMockTransport()
		p3 := httputil.NewSingleHostReverseProxy(&url.URL{})
		p3.Transport = tr3
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		rb := rebalancer.New([]*httputil.ReverseProxy{p1, p2, p3}, mockLb, logger)
		go rb.StartRebalancer(time.Millisecond*200, "https://foo\\x7f.com/")
		time.Sleep(time.Microsecond * 250)
	})
	t.Run("it returns the same load balancer that was provided", func(t *testing.T) {
		lb := roundrobin.New[*httputil.ReverseProxy](nil)
		rb := rebalancer.New(nil, lb, nil)
		assert.Same(t, lb, rb.GetLoadbalancer())
	})
}
