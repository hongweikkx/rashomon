package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/router/handle"
	"github.com/hongweikkx/rashomon/router/middleware/auth"
	"github.com/hongweikkx/rashomon/router/middleware/ctxkv"
	"github.com/hongweikkx/rashomon/router/middleware/mlog"
	"github.com/hongweikkx/rashomon/router/middleware/rate"
	"github.com/hongweikkx/rashomon/service"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
)

func Router(engine *gin.Engine) {
	engine.Use(
		ctxkv.Bind,
		mlog.LogMiddle,
		mlog.RecoveryWithLog,
		nrgin.Middleware(service.MonitorNewRelic),
		gzip.Gzip(gzip.DefaultCompression),
	)

	engine.NoRoute(handle.NoRoute)
	engine.GET("/", handle.Health)
	engine.GET("/health", handle.Health)

	v1 := engine.Group("v1")
	userR := v1.Group("user")
	userR.POST("login", rate.LimitDefault(), modAuth.AuthIns.LoginHandler)
	userR.POST("logout", rate.LimitDefault(), modAuth.AuthIns.MiddlewareFunc(), modAuth.AuthIns.LogoutHandler)
	userR.POST("refresh_token", rate.LimitDefault(), modAuth.AuthIns.MiddlewareFunc(), modAuth.AuthIns.RefreshHandler)
	userR.GET("info", rate.LimitDefault(), modAuth.AuthIns.MiddlewareFunc(), handle.UserInfo)

	answerR := v1.Group("answer", modAuth.AuthIns.MiddlewareFunc())
	answerR.GET("recommend", rate.LimitDefault(), handle.AnswerList)
	answerR.GET("likelist", rate.LimitDefault(), handle.LikeList)
	answerR.POST("like", rate.LimitDefault(), handle.Like)
	answerR.POST("search", rate.LimitDefault(), handle.Search)

	toolR := v1.Group("tool", modAuth.AuthIns.MiddlewareFunc())
	toolR.GET("list", rate.LimitDefault(), handle.ToolList)
}
