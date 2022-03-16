//
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/router"
)

func main() {
	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	router.Router(engine)
	serv := &http.Server{
		Addr:    conf.AppConfig.Addr,
		Handler: engine,
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Fatal("http serv err:", err.Error())
		}
	}()
	log.SugarLogger.Infof("server started on %+v...", conf.AppConfig.Addr)
	// wait to stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		log.SugarLogger.Info("[dasboard] http server forced to shutdown:", err)
	}
	log.SugarLogger.Info("server exit.")
	log.SugarLogger.Sync()
}
