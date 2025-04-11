package log

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/helper"
)

func NewLogger(appConfig *config.AppConfig) *logrus.Logger {
	var formatter logrus.Formatter

	if appConfig.LogFormatter == "text" {
		formatter = &logrus.TextFormatter{FullTimestamp: true}
	} else {
		formatter = &logrus.JSONFormatter{}
	}

	logLevel, err := logrus.ParseLevel(appConfig.LogLevel)
	if err != nil {
		panic(err)
	}

	return &logrus.Logger{
		Out:          os.Stdout,
		Formatter:    formatter,
		ReportCaller: helper.ToBool(appConfig.LogReportCaller),
		Level:        logLevel,
	}
}
