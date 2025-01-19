package main

import (
	"fmt"

	"github.com/cyberhck/roundguard/config"
	"github.com/cyberhck/roundguard/pkg/http/echo"
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

	err = cmd.Execute()
	if err != nil {
		logrus.WithField("error", err).Panic("failed to execute command")
	}
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
