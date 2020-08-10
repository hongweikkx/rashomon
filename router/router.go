package router

import (
	"github.com/afex/hystrix-go/hystrix"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"net/http"
	"time"
)

func Use(engine *gin.Engine) {
	engine.Use(
		ginzap.Ginzap(log.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(log.Logger, true),
	)

	engine.GET("/pingDegrade", handlePingDegrade)
	engine.GET("/pingFuse", handlePingFuse)
}

func handlePingFuse(c *gin.Context) {
	hystrix.Go("fuse", func() error {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		return nil
	}, func(err error) error {
		c.String(http.StatusInternalServerError, "service degrade")
		return nil
	})
}

func handlePingDegrade(c *gin.Context) {
	hystrix.Go("degrade", func() error {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		return nil
	}, func(err error) error {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status": "degrade",
		})
		return nil
	})
}