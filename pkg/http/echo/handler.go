package echo

import (
	"encoding/json"
	"github.com/cyberhck/roundguard/pkg/http/types"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Handler(logger *logrus.Logger) http.HandlerFunc {
	errBadRequest := types.ErrorResponse{
		Title:       "Bad Request",
		Description: "The request body could not be parsed as JSON",
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var data any
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.WithField("error", err).Warn("failed parsing request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(errBadRequest)
			if err != nil {
				logger.WithField("error", err).Error("failed writing request")
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.WithField("error", err).Error("failed writing request")
		}
	}
}
