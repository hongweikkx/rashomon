package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"rashomon/api/response"
	"rashomon/middleware/auth"
	"rashomon/middleware/ctxkv"
	"rashomon/middleware/mlog"
	"rashomon/middleware/rate"
	"rashomon/router/handle"
)

func Router(engine *gin.Engine) {
	engine.Use(
		ctxkv.Bind,
		mlog.LogMiddle,
		mlog.RecoveryWithLog,
		gzip.Gzip(gzip.DefaultCompression),
	)

	engine.NoRoute(response.NoRoute)
	engine.GET("/", response.Health)
	engine.GET("/health", response.Health)

	v1 := engine.Group("v1")
	authG := v1.Group("/")
	authG.Use(auth.Auth.MiddlewareFunc())
	v1.POST("login", rate.LimitDefault(), auth.Auth.LoginHandler)
	authG.POST("logout", rate.LimitDefault(), auth.Auth.LogoutHandler)
	authG.POST("refresh_token", rate.LimitDefault(), auth.Auth.RefreshHandler)
	authG.GET("user", rate.LimitDefault(), handle.UserInfo)
}
