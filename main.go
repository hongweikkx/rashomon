//
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/etcd"
	"github.com/hongweikkx/rashomon/hystrix"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/router"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// log
	log.InitLogger()
	defer log.Logger.Sync()
	defer log.SugarLogger.Sync()

	// conf
	err := conf.InitConf()
	if err != nil {
		log.SugarLogger.Fatal("conf err:", err.Error())
		return
	}
	// etcd
	if conf.AppConfig.ETCD.Enable {
		// todo 应该是有多个这种  因为api server 接的是多组 而非一组
		master, err := etcd.New("master/", []etcd.EndPoint{})
		if err != nil {
			log.SugarLogger.Fatal("etcd err:", err.Error())
			return
		}
		defer master.Cli.Close()
	}
	// hystrix
	hystrix.InitHystrix()

	// http server
	engine := gin.New()
	router.Use(engine)
	serv := &http.Server{
		Addr: conf.AppConfig.HttpServer.Addr,
		Handler: engine,
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Error("serv err:", err.Error())
		}
	}()
	// sig -> shutdown server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.SugarLogger.Error("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		log.SugarLogger.Error("Server forced to shutdown:", err)
	}
	log.SugarLogger.Error("Server exiting")
}