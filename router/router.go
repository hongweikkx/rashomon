package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"rashomon/api/controller"
	"rashomon/api/response"
	"rashomon/middleware/auth"
	"rashomon/middleware/ctxkv"
	"rashomon/middleware/mlog"
)

func Router(engine *gin.Engine) {
	engine.Use(
		ctxkv.Bind,
		mlog.LogMiddle,
		mlog.RecoveryWithLog,
		auth.MiddleWare(auth.Auth),
		gzip.Gzip(gzip.DefaultCompression),
	)

	engine.NoRoute(response.NoRoute)
	engine.GET("/", response.Health)
	engine.GET("/health", response.Health)

	authG := engine.Group("/")
	authG.Use(auth.Auth.MiddlewareFunc())
	engine.POST("user/login", auth.Auth.LoginHandler)
	authG.POST("user/logout", auth.Auth.LogoutHandler)
	authG.POST("user/refresh_token", auth.Auth.RefreshHandler)
	userC := controller.NewUserController()
	authG.GET("user/info", userC.UserInfo)
}
