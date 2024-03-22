package log

import (
	"fmt"
	"os"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sirupsen/logrus"
)

func NewLogger(appConfig *config.AppConfig) *logrus.Logger {

	var formatter logrus.Formatter

	if appConfig.Log_Formatter == "text" {
		formatter = &logrus.TextFormatter{FullTimestamp: true}
	} else {
		formatter = &logrus.JSONFormatter{}
	}

	logLevel, err := logrus.ParseLevel(appConfig.Log_Level)
	if err != nil {
		fmt.Println(appConfig)
		panic(err)
	}

	return &logrus.Logger{
		Out:          os.Stdout,
		Formatter:    formatter,
		ReportCaller: helper.ToBool(appConfig.Log_ReportCaller),
		Level:        logLevel,
	}
}
