package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/middleware/auth"
	"github.com/hongweikkx/rashomon/middleware/hystrix"
	"net/http"
	"time"
)

func Use(engine *gin.Engine) {
	authMiddleware := auth.New()
	auth.Use(authMiddleware, engine)
	engine.NoRoute(handleNoRoute)
	//middleware
	engine.Use(
		ginzap.Ginzap(log.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(log.Logger, true),
	)
	engine.GET("/pingDegrade", hystrix.HandleFuse, handlePing)
	engine.GET("/pingFuse", hystrix.HandleFuse, handlePing)
}


func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}


func handleNoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}