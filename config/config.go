package config

import (
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	LogLevel  logrus.Level `env:"LOG_LEVEL" default:"info"`
	LogFormat string       `env:"LOG_FORMAT" default:"text"` // or use JSON
}

func Load[T any]() (*T, error) {
	var cfg T
	err := configor.Load(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
