package router

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/router/handle"
	"github.com/hongweikkx/rashomon/router/middleware/auth"
)

func Router(engine *gin.Engine) {
	engine.NoRoute(handle.NoRoute)
	engine.GET("/", handle.Health)
	engine.GET("/health", handle.Health)
	//middleware
	engine.Use(
		ginzap.Ginzap(log.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(log.Logger, true),
	)

	v1 := engine.Group("v1")

	userR := v1.Group("user")
	userR.POST("login", auth.AuthMiddleWare.LoginHandler)
	userR.POST("logout", auth.AuthMiddleWare.MiddlewareFunc(), auth.AuthMiddleWare.LogoutHandler)
	userR.POST("refresh_token", auth.AuthMiddleWare.MiddlewareFunc(), auth.AuthMiddleWare.RefreshHandler)
	userR.GET("info", auth.AuthMiddleWare.MiddlewareFunc(), handle.UserInfo)

	answerR := v1.Group("answer", auth.AuthMiddleWare.MiddlewareFunc())
	answerR.GET("recommend", handle.AnswerList)
	answerR.GET("likelist", handle.LikeList)
	answerR.POST("like", handle.Like)
	answerR.POST("search", handle.Search)

	toolR := v1.Group("tool", auth.AuthMiddleWare.MiddlewareFunc())
	toolR.GET("list", handle.ToolList)
}
