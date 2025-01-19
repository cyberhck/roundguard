package config_test

import (
	"testing"

	"github.com/cyberhck/roundguard/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type LogLevelTestCase struct {
	input    string
	expected logrus.Level
	name     string
}

func TestLoad(t *testing.T) {
	testCases := []LogLevelTestCase{
		{input: "debug", expected: logrus.DebugLevel},
		{input: "info", expected: logrus.InfoLevel},
		{input: "warn", expected: logrus.WarnLevel},
		{input: "error", expected: logrus.ErrorLevel},
		{input: "fatal", expected: logrus.FatalLevel},
		{input: "panic", expected: logrus.PanicLevel},
		{input: "trace", expected: logrus.TraceLevel},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tc.input)
			cfg, err := config.Load[config.AppConfig]()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, cfg.LogLevel)
		})
	}
}
