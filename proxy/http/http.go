package proxyhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/router"
	"context"
	"net/http"
	"time"
)

func Start() *http.Server{
	engine := gin.New()
	router.Use(engine)
	serv := &http.Server{
		Addr: conf.AppConfig.Proxy.HttpServer.Addr,
		Handler: engine,
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Error("http serv err:", err.Error())
		}
	}()
	return serv
}

func Stop(httpServer *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.SugarLogger.Error("http server forced to shutdown:", err)
	}
	log.SugarLogger.Error("http server exit")
}
