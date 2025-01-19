package echo

import (
	"net/http"

	"github.com/cyberhck/roundguard/pkg/http/echo/reflect"
	"github.com/sirupsen/logrus"

	"github.com/cyberhck/roundguard/pkg/http/echo/liveness"
)

func CreateServer(address string, logger *logrus.Logger) *http.Server {
	router := http.NewServeMux()
	router.Handle("GET /live", liveness.Handler())
	router.Handle("POST /reflect", reflect.Handler(logger))

	return &http.Server{Addr: address, Handler: router}
}
