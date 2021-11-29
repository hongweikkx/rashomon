package proxy

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/dashboard/route/handle"
	"github.com/hongweikkx/rashomon/log"
	"time"
)

func Router(engine *gin.Engine) error {
	engine.NoRoute(handle.NoRoute)
	engine.GET("/", handle.Health)
	engine.GET("/health", handle.Health)
	//middleware
	engine.Use(
		ginzap.Ginzap(log.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(log.Logger, true),
	)
	return nil
}