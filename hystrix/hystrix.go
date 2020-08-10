package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hongweikkx/rashomon/conf"
)

func InitHystrix() {
	hystrix.ConfigureCommand("degrade", hystrix.CommandConfig{
		Timeout:               conf.AppConfig.Hystrix.Degrade.Timeout,
		MaxConcurrentRequests: conf.AppConfig.Hystrix.Degrade.MaxConcurrentRequests,
		RequestVolumeThreshold: conf.AppConfig.Hystrix.Degrade.RequestVolumeThreshold,
		ErrorPercentThreshold: conf.AppConfig.Hystrix.Degrade.ErrorPercentThreshold,
	})
	hystrix.ConfigureCommand("fuse", hystrix.CommandConfig{
		Timeout:               conf.AppConfig.Hystrix.Fuse.Timeout,
		MaxConcurrentRequests: conf.AppConfig.Hystrix.Fuse.MaxConcurrentRequests,
		RequestVolumeThreshold: conf.AppConfig.Hystrix.Fuse.RequestVolumeThreshold,
		ErrorPercentThreshold: conf.AppConfig.Hystrix.Fuse.ErrorPercentThreshold,
	})
}
