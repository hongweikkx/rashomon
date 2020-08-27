//
package main

import (
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/dashboard"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/proxy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Init()
	err := conf.Init()
	if err != nil {
		log.SugarLogger.Fatal("conf err:", err.Error())
		return
	}
	// proxy
	proxy.Start()
	// dashboard
	dashboard.Start()
	log.SugarLogger.Info("server started...")
	// wait to stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	proxy.Stop()
	dashboard.Stop()
	log.SugarLogger.Info("server exit.")
	log.Stop()
}