package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hongweikkx/rashomon/conf"
	"go.etcd.io/etcd/clientv3"
	"github.com/hongweikkx/rashomon/log"
	"time"
)



type Service struct {
	Name    string
	Info    ServiceInfo
	Cli     *clientv3.Client
	Stop    chan error
	LeaseId clientv3.LeaseID
}

var ServiceRegClient Service

func ServiceRegInitClient(name string, info ServiceInfo) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: conf.AppConfig.ETCD.EndPoints,
		DialTimeout: time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second,
	})
	if err != nil {
		log.SugarLogger.Error("service reg client error:", err.Error())
	}
	ServiceRegClient = Service{Name: name, Info: info, Cli: cli, Stop: make(chan error)}
	err = ServiceRegClient.Start()
	if err != nil {
		log.SugarLogger.Error("service stop with error:", err.Error())
	}
}


func (service *Service) Start() error{
	keepAliveCh, err := service.KeepAlive()
	if err != nil {
		log.SugarLogger.Error("etcd error:", err.Error())
		return err
	}
	for {
		select {
		case errStop := <-service.Stop:
			service.Revoke()
			return errStop
		case <- service.Cli.Ctx().Done():
			return errors.New("server close")
		case re, ok := <- keepAliveCh:
			if !ok {
				service.Revoke()
				return errors.New("server keepalive close")
			} else {
				fmt.Print("", re)
			}
		}
	}
}

func (service *Service) KeepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error){
	// 就是监控key
	key := conf.AppConfig.ETCD.WatchPrix + service.Name
	value, _ := json.Marshal(service.Info)
	resp, err := service.Cli.Grant(context.TODO(), 5)
	if err != nil {
		return nil, err
	}
	_, err = service.Cli.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}
	service.LeaseId = resp.ID
	return service.Cli.KeepAlive(context.TODO(), resp.ID)
}


func (service *Service)Revoke() error{
	_, err := service.Cli.Revoke(context.TODO(), service.LeaseId)
	if err != nil {
		log.SugarLogger.Error("err:", err.Error())
	}
	return nil
}

func (service *Service) StopService()  {
	service.Stop <- nil
}


