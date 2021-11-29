package proxy

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/load_balance"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/storage"
	"net/http"
	"time"
)

type Proxy struct {
	Clusters   []*load_balance.Cluster
	StoreCli   storage.Storeage
	HttpServer *http.Server
}

var ProxyIns *Proxy

func Start() error {
	httpServer, err := ProxyIns.StartHttp()
	if err != nil {
		return err
	}
	ProxyIns = &Proxy{
		Clusters:   nil,
		HttpServer: httpServer,
	}
	return nil
}

func Stop() {
	ProxyIns.StopHttp()
	ProxyIns.StoreCli.Stop()
}

func (proxy *Proxy) StartHttp() (*http.Server, error) {
	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	err := Router(engine)
	if err !=nil {
		return nil, err
	}
	serv := &http.Server{
		Addr:    conf.AppConfig.Proxy.HTTPServer.Addr,
		Handler: engine,
	}
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Fatal("[proxy] http serv err:", err.Error())
		}
	}()
	return serv, nil
}

func (proxy *Proxy) StopHttp() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := proxy.HttpServer.Shutdown(ctx); err != nil {
		log.SugarLogger.Info("[proxy] http server forced to shutdown:", err)
	}
}