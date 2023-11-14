package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
)

func GinLog(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		start := time.Now()

		c.Next()

		stop := time.Since(start)
		responseTime := float32(stop.Nanoseconds()) / 1000000.0
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		requestUri := c.Request.RequestURI

		matchKubeProbeUserAgent, _ := regexp.MatchString("^(kube-probe.*)$", clientUserAgent)

		if clientUserAgent == "Envoy/HC" || clientUserAgent == "Consul Health Check" || matchKubeProbeUserAgent {
			return
		}

		fields := logrus.Fields{
			"statusCode":     statusCode,
			"responseTime":   fmt.Sprintf("%.2fms", responseTime),
			"clientIP":       clientIP,
			"requestMethod":  c.Request.Method,
			"requestPath":    path,
			"requestReferer": referer,
			"userAgent":      clientUserAgent,
			"requestUri":     requestUri,
		}

		if len(c.Errors) > 0 {
			logger.WithFields(fields).Errorf(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode > 499 {
				logger.WithFields(fields).Errorf("HTTP Status Failed")
			} else if statusCode > 399 {
				logger.WithFields(fields).Warnf("HTTP Status Failed")
			} else {
				logger.WithFields(fields).Infof("HTTP Status OK")
			}
		}
	}
}
