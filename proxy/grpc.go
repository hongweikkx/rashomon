package proxy

import (
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"google.golang.org/grpc"
	"net"
)

func (proxy *Proxy) StartGRPC() *grpc.Server {
	lis, err := net.Listen("tcp", conf.AppConfig.Proxy.GrpcServer.Addr)
	if err != nil {
		panic("grpc server panic:" + err.Error())
	}
	s := grpc.NewServer()
	go func() {
		if err := s.Serve(lis); err != nil {
			log.SugarLogger.Fatal("grpc serv err:", err.Error())
		}
	}()
	return s
}

func (proxy *Proxy) StopGRPC() {
	proxy.GrpcServer.Stop()
}
