package echo_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyberhck/roundguard/pkg/http/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	t.Run("calling /live endpoint should return status 200", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://localhost:5555/live", nil)
		assert.NoError(t, err)
		echo.CreateServer("http://localhost:8080", logger).Handler.ServeHTTP(responseWriter, req)
		assert.Equal(t, http.StatusOK, responseWriter.Code)
	})
	t.Run("calling /reflect endpoint should return 400 if the request body nil", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "http://localhost:5555/reflect", nil)
		assert.NoError(t, err)
		echo.CreateServer("http://localhost:8080", logger).Handler.ServeHTTP(responseWriter, req)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
	})
	t.Run("calling /reflect endpoint should return 400 if the request body is not json", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		const mockInvalidJson = "{something>"
		req, err := http.NewRequest("POST", "http://localhost:5555/reflect", bytes.NewBufferString(mockInvalidJson))
		assert.NoError(t, err)
		echo.CreateServer("http://localhost:8080", logger).Handler.ServeHTTP(responseWriter, req)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
	})
	t.Run("calling /reflect endpoint should return the same value with 200 status if the body is valid JSON", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		const validJson = `{"hello":"world"}`
		req, err := http.NewRequest("POST", "http://localhost:5555/reflect", bytes.NewBufferString(validJson))
		assert.NoError(t, err)
		echo.CreateServer("http://localhost:8080", logger).Handler.ServeHTTP(responseWriter, req)
		assert.Equal(t, http.StatusOK, responseWriter.Code)
		assert.Equal(t, validJson, strings.Trim(responseWriter.Body.String(), "\n")) // json.Encode adds a new line as described here: https://github.com/golang/go/blob/40b3c0e58a0ae8dec4684a009bf3806769e0fc41/src/encoding/json/stream.go#L215
	})
}
