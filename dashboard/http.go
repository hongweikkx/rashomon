package dashboard

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/middleware/auth"
	"net/http"
	"time"
)

type Dashboard struct {
	Serv *http.Server
}

var DashboardIns *Dashboard

func Start() error{
	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	err := router(engine)
	if err != nil {
		return err
	}
	serv := &http.Server{
		Addr: conf.AppConfig.DashBoard.Addr,
		Handler: engine,
	}
	DashboardIns = &Dashboard{Serv: serv}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Fatal("http serv err:", err.Error())
		}
	}()
	return nil
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := DashboardIns.Serv.Shutdown(ctx); err != nil {
		log.SugarLogger.Info("http server forced to shutdown:", err)
	}
}

func router(engine *gin.Engine) error{
	authMiddleware, err := auth.New()
	if err != nil {
		return err
	}
	auth.Use(authMiddleware, engine)
	engine.NoRoute(handleNoRoute)
	//middleware
	engine.Use(
		ginzap.Ginzap(log.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(log.Logger, true),
	)
	engine.GET("/ping", handlePing)
	return nil
}


func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handleNoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}

