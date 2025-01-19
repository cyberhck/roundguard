package reflect

import (
	"encoding/json"
	"net/http"

	"github.com/cyberhck/roundguard/pkg/http/types"
	"github.com/sirupsen/logrus"
)

func Handler(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(types.ErrorResponse{
				Title:       "Empty request body",
				Description: "A valid JSON body is required for this operation",
			})
			if err != nil {
				logger.WithField("error", err).Error("Failed writing response")
			}
			return
		}
		var data any
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.WithField("error", err).Warn("failed parsing request body")
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(types.ErrorResponse{
				Title:       "Bad Request",
				Description: "The request body could not be parsed as JSON",
			})
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
