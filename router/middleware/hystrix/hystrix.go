package hystrix

import (
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
)

func init() {
	hystrix.ConfigureCommand("degrade", hystrix.CommandConfig{
		Timeout:                conf.AppConfig.Hystrix.Degrade.Timeout,
		MaxConcurrentRequests:  conf.AppConfig.Hystrix.Degrade.MaxConcurrentRequests,
		RequestVolumeThreshold: conf.AppConfig.Hystrix.Degrade.RequestVolumeThreshold,
		ErrorPercentThreshold:  conf.AppConfig.Hystrix.Degrade.ErrorPercentThreshold,
	})
	hystrix.ConfigureCommand("fuse", hystrix.CommandConfig{
		Timeout:                conf.AppConfig.Hystrix.Fuse.Timeout,
		MaxConcurrentRequests:  conf.AppConfig.Hystrix.Fuse.MaxConcurrentRequests,
		RequestVolumeThreshold: conf.AppConfig.Hystrix.Fuse.RequestVolumeThreshold,
		ErrorPercentThreshold:  conf.AppConfig.Hystrix.Fuse.ErrorPercentThreshold,
	})
}

func HandleFuse(c *gin.Context) {
	hystrix.Go("fuse", func() error {
		c.Next()
		return nil
	}, func(err error) error {
		c.String(http.StatusInternalServerError, "service degrade")
		return nil
	})
}

func HandleDegrade(c *gin.Context) {
	hystrix.Go("degrade", func() error {
		c.Next()
		return nil
	}, func(err error) error {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "service degrade",
		})
		return nil
	})
}
