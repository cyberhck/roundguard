package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/cyberhck/roundguard/config"
	"github.com/cyberhck/roundguard/pkg/http/echo"
	"github.com/cyberhck/roundguard/pkg/http/loadbalancer"
	"github.com/cyberhck/roundguard/pkg/http/transports"
	"github.com/cyberhck/roundguard/pkg/rebalancer"
	"github.com/cyberhck/roundguard/pkg/roundrobin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	cfg, err := config.Load[config.AppConfig]()
	if err != nil {
		logrus.Fatal(err)
	}
	logger := createLogger(cfg)
	cmd := &cobra.Command{
		Use:   "roundguard",
		Short: "command line application that starts different servers",
		Long:  `A reverse proxy with load balancing capability coupled with a JSON echo server`,
	}
	setupEchoServerCommand(cmd, logger)
	setupLoadBalancerCommand(cmd, logger, cfg)

	err = cmd.Execute()
	if err != nil {
		logrus.WithField("error", err).Panic("failed to execute command")
	}
}

func setupLoadBalancerCommand(cmd *cobra.Command, logger *logrus.Logger, cfg *config.AppConfig) {
	lbCommand := &cobra.Command{
		Use:   "lb",
		Short: "manage load balancing server",
	}
	startLbCommand := &cobra.Command{
		Use:   "start",
		Short: "start load balancer",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				return err
			}
			hosts, err := cmd.Flags().GetStringArray("hosts")
			if err != nil {
				return err
			}
			logger.WithField("hosts", hosts).Info("starting load balancer with hosts")
			targets := lo.Map(hosts, func(item string, _ int) *httputil.ReverseProxy {
				result, err := url.Parse(item)
				if err != nil {
					panic(err)
				}

				proxy := httputil.NewSingleHostReverseProxy(result)
				proxy.Transport = transports.NewMeasurement(http.DefaultTransport)

				return proxy
			})
			rb := rebalancer.New(targets, roundrobin.New(targets), logger.WithField("scope", "rebalancer").Logger)
			go rb.StartRebalancer(cfg.ProxyCheckDuration, "/live")

			address := fmt.Sprintf(":%d", port)
			server, err := loadbalancer.CreateServer(address, rb, logger.WithField("scope", "load_balancer").Logger)
			if err != nil {
				return err
			}
			logger.WithField("address", fmt.Sprintf("http://localhost%s", address)).Info("starting load balancer")

			return server.ListenAndServe()
		},
	}
	startLbCommand.Flags().Int("port", 3000, "port to start")
	startLbCommand.Flags().StringArray("hosts", []string{"localhost:8000", "localhost:8001", "localhost:8002"}, "hosts to start")
	lbCommand.AddCommand(startLbCommand)
	cmd.AddCommand(lbCommand)
}

func createLogger(cfg *config.AppConfig) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(cfg.LogLevel)
	logger.SetFormatter(&logrus.TextFormatter{})
	switch cfg.LogFormat {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
		return logger
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	return logger.WithField("scope", "startup").Logger
}

func setupEchoServerCommand(cmd *cobra.Command, logger *logrus.Logger) {
	echoCommand := &cobra.Command{
		Use:   "echo",
		Short: "manage echo server",
	}
	startEchoCommand := &cobra.Command{
		Use:   "start",
		Short: "start echo server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				return err
			}
			address := fmt.Sprintf(":%d", port)
			server := echo.CreateServer(address, logger.WithField("scope", "echo_server").Logger)
			logger.WithField("address", fmt.Sprintf("http://localhost%s", address)).Info("starting echo server")

			return server.ListenAndServe()
		},
	}
	startEchoCommand.Flags().Int("port", 8000, "port to start")
	echoCommand.AddCommand(startEchoCommand)
	cmd.AddCommand(echoCommand)
}
