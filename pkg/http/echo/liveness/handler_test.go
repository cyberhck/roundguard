package liveness_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyberhck/roundguard/pkg/http/echo/liveness"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("returns 200 OK status", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		liveness.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
