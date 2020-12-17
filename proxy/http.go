package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"github.com/valyala/fasthttp"
	"net/http"
)

func (proxy *Proxy) StartHttp() *fasthttp.Server {
	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	serv := &fasthttp.Server{
		Handler: proxy.SrvHTTP,
	}
	go func() {
		if err := serv.ListenAndServe(conf.AppConfig.Proxy.HttpServer.Addr); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Fatal("http serv err:", err.Error())
		}
	}()
	return serv
}

func (proxy *Proxy) StopHttp() {
	if err := proxy.HttpServer.Shutdown(); err != nil {
		log.SugarLogger.Info("http server forced to shutdown:", err)
	}
}

func (p *Proxy) SrvHTTP(ctx *fasthttp.RequestCtx) {
	// 分析ctx的url
	//

}
