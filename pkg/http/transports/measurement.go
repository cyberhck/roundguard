package transports

import "net/http"

type measurement struct {
	roundTripper http.RoundTripper
}

func (m *measurement) RoundTrip(request *http.Request) (*http.Response, error) {
	res, err := m.roundTripper.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	res.Header.Add("X-Backend-Server", request.URL.Host)

	return res, nil
}

func NewMeasurement(roundTripper http.RoundTripper) http.RoundTripper {
	return &measurement{roundTripper}
}
