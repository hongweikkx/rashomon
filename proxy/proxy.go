package proxy

import (
	"github.com/hongweikkx/rashomon/load_balance"
	"github.com/hongweikkx/rashomon/log"
	proxygrpc "github.com/hongweikkx/rashomon/proxy/grpc"
	proxyhttp "github.com/hongweikkx/rashomon/proxy/http"
	"github.com/hongweikkx/rashomon/storage"
	"google.golang.org/grpc"
	"net/http"
)

type Proxy struct {
	Clusters []*load_balance.Cluster
	StoreCli storage.Storeage
	HttpServer *http.Server
	GrpcServer *grpc.Server
}


var ProxyIns *Proxy

func Start() {
	httpServer := proxyhttp.Start()
	grpcServer := proxygrpc.Start()
	storeCli, err := storage.Start()
	if err != nil {
		log.SugarLogger.Fatal("etcd error:", err.Error())
	}
	ProxyIns = &Proxy{
		Clusters: nil,
		StoreCli: storeCli,
		HttpServer: httpServer,
		GrpcServer: grpcServer,
	}
}

func Stop() {
	proxyhttp.Stop(ProxyIns.HttpServer)
	proxygrpc.Stop(ProxyIns.GrpcServer)
	ProxyIns.StoreCli.Stop()
}







