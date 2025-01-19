package liveness

import "net/http"

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // we'll hardcode the liveness to return 200 OK.
		// todo: replace with https://github.com/hellofresh/health-go
	}
}
