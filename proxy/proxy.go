package proxy

import (
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/storage"
	"os"
	"os/signal"
	"syscall"
	proxygrpc "github.com/hongweikkx/rashomon/proxy/grpc"
	proxyhttp "github.com/hongweikkx/rashomon/proxy/http"
)

type Proxy struct {
	Clusters []*Cluster
	StoreCli storage.Storeage
}

type Cluster struct {
	Servers []*Server
	LBStrategy int
}

type Server struct {
	Addr string
	Weight int
}

var ProxyIns *Proxy

func Start() {
	httpServer := proxyhttp.Start()
	GrpcServer := proxygrpc.Start()
	storeCli, err := storage.StartStorage()
	if err != nil {
		log.SugarLogger.Fatal("etcd error:", err.Error())
	}
	ProxyIns = &Proxy{
		Clusters: nil,
		StoreCli: storeCli,
	}
	// stop k
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.SugarLogger.Error("Shutting down server...")
	proxyhttp.Stop(httpServer)
	proxygrpc.Stop(GrpcServer)
	ProxyIns.StoreCli.Stop()
	log.SugarLogger.Error("server exit.")
}







