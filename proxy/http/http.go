package proxyhttp

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/middleware/hystrix"
	"context"
	"net/http"
	"time"
)

func Start() *http.Server{
	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	router(engine)
	serv := &http.Server{
		Addr: conf.AppConfig.Proxy.HttpServer.Addr,
		Handler: engine,
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Fatal("http serv err:", err.Error())
		}
	}()
	return serv
}

func Stop(httpServer *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.SugarLogger.Info("http server forced to shutdown:", err)
	}
}

func router(engine *gin.Engine) {
	engine.NoRoute(handleNoRoute)
	//middleware
	engine.Use(
		ginzap.Ginzap(log.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(log.Logger, true),
	)
	engine.GET("/pingDegrade", hystrix.HandleFuse, handlePing)
	engine.GET("/pingFuse", hystrix.HandleFuse, handlePing)
}

func handleNoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}


func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
