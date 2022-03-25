package mlog

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/router/middleware/ctxkv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogMiddle(c *gin.Context) {
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()
	logger := ctxkv.Log(c)
	end := time.Now()
	latency := end.Sub(start)
	end = end.UTC()

	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("time", end.Format(time.RFC3339)),
		}
		if ctxkv.GetDgd(c) {
			fields = append(fields, zap.Bool("degrade", true))
		}
		logger.Info(path, fields...)
	}
}

func RecoveryWithLog(c *gin.Context) {
	logger := ctxkv.Log(c)
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			if brokenPipe {
				logger.Error(c.Request.URL.Path,
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
				)
				// If the connection is dead, we can't write a status to it.
				c.Error(err.(error)) // nolint: errcheck
				c.Abort()
				return
			}

			logger.Error("[RECOVER FROM PANIC]",
				zap.Any("error", err),
				zap.String("time", time.Now().Format(time.RFC3339)),
			)
			c.String(http.StatusInternalServerError, "")
		}
	}()
	c.Next()
}
