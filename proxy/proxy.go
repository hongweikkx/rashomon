package proxy

import (
	"github.com/hongweikkx/rashomon/load_balance"
	"github.com/hongweikkx/rashomon/storage"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
)

type Proxy struct {
	Clusters   []*load_balance.Cluster
	StoreCli   storage.Storeage
	HttpServer *fasthttp.Server
	GrpcServer *grpc.Server
}

var ProxyIns *Proxy

func Start() error {
	storeCli, err := storage.Start()
	httpServer := ProxyIns.StartHttp()
	grpcServer := ProxyIns.StartGRPC()
	if err != nil {
		return err
	}
	ProxyIns = &Proxy{
		Clusters:   nil,
		StoreCli:   storeCli,
		HttpServer: httpServer,
		GrpcServer: grpcServer,
	}
	return nil
}

func Stop() {
	ProxyIns.StopGRPC()
	ProxyIns.StopHttp()
	ProxyIns.StoreCli.Stop()
}
