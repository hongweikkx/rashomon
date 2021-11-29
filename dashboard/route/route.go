package route

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/dashboard/route/handle"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/middleware/auth"
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

	v1 := engine.Group("/v1/dashboard")
	v1.POST("/user/login", auth.AuthMiddleWare.LoginHandler)
	v1.POST("/user/logout",auth.AuthMiddleWare.MiddlewareFunc(), auth.AuthMiddleWare.LogoutHandler)
	v1.POST("/user/refresh_token", auth.AuthMiddleWare.MiddlewareFunc(), auth.AuthMiddleWare.RefreshHandler)
	v1.GET("/user/info", auth.AuthMiddleWare.MiddlewareFunc(), handle.UserInfo)
	v1.GET("/table/list", auth.AuthMiddleWare.MiddlewareFunc(), handle.TableList)
	return nil
}

