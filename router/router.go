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

	authG := engine.Group("/")
	authG.Use(auth.Auth.MiddlewareFunc())
	engine.POST("user/login", rate.LimitDefault(), auth.Auth.LoginHandler)
	authG.POST("user/logout", rate.LimitDefault(), auth.Auth.LogoutHandler)
	authG.POST("user/refresh_token", rate.LimitDefault(), auth.Auth.RefreshHandler)
	authG.GET("user/info", rate.LimitDefault(), handle.UserInfo)
}
