//
package main

import (
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/proxy"
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
	// proxy
	proxy.Start()
}