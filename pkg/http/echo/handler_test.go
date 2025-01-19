package echo_test

import (
	"bytes"
	"github.com/cyberhck/roundguard/pkg/http/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	t.Run("if body is not valid JSON, it returns a status 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString("invalid json"))
		echo.Handler(logger).ServeHTTP(w, request)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("if body is valid JSON, it emits it back", func(t *testing.T) {
		const mockJson = `{"hello":"world"}`
		request := httptest.NewRequest("POST", "/", bytes.NewBufferString(mockJson))
		w := httptest.NewRecorder()
		echo.Handler(logger).ServeHTTP(w, request)
		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := strings.Trim(w.Body.String(), "\n") // json.Encode adds a new line as described here: https://github.com/golang/go/blob/40b3c0e58a0ae8dec4684a009bf3806769e0fc41/src/encoding/json/stream.go#L215
		assert.Equal(t, mockJson, responseBody)
	})
}
