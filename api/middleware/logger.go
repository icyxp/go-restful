package middleware

import (
	"fmt"
	"go-restful/lib/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	timeFormat = "02/Jan/2006:15:04:05 -0700"
)

//LoggerInit 实例化
func LoggerInit() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	return logger
}

//Logger 格式化日志输出
func Logger() gin.HandlerFunc {
	logger := LoggerInit()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	return func(c *gin.Context) {
		// other handler can change c.Path so:
		header := c.Request.Header
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		//latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		traceID := c.Request.Header.Get("x-amzn-trace-id")
		requestID, _ := c.Get(utils.XRequestID)
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"hostname":   fmt.Sprintf("%s", hostname),
			"statusCode": statusCode,
			"latency":    fmt.Sprintf("%s", latency),
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"dataLength": dataLength,
			"traceID":    traceID,
			"requestID":  requestID,
			"header":     header,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, traceID, requestID, latency)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
