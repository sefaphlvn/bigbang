package log

import (
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger(appConfig *config.AppConfig) *logrus.Logger {

	var formatter logrus.Formatter

	if appConfig.Log.Formatter == "text" {
		formatter = &logrus.TextFormatter{FullTimestamp: true}
	} else {
		formatter = &logrus.JSONFormatter{}
	}

	logLevel, err := logrus.ParseLevel(appConfig.Log.Level)
	if err != nil {
		panic(err)
	}

	return &logrus.Logger{
		Out:          os.Stdout,
		Formatter:    formatter,
		ReportCaller: appConfig.Log.ReportCaller,
		Level:        logLevel,
	}
}
