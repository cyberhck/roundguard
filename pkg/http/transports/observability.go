package transports

import (
	"net/http"

	"github.com/cyberhck/roundguard/pkg/measurements"
)

type observability struct {
	roundTripper http.RoundTripper
}

func (m *observability) RoundTrip(request *http.Request) (*http.Response, error) {
	timer := measurements.MeasureDuration("proxy", request.Method, request.URL.Host)
	res, err := m.roundTripper.RoundTrip(request)
	if err != nil {
		timer("fail")

		return nil, err
	}
	res.Header.Add("X-Backend-Server", request.URL.Host)
	timer("success")

	return res, nil
}

func NewMeasurement(roundTripper http.RoundTripper) http.RoundTripper {
	return &observability{roundTripper}
}
