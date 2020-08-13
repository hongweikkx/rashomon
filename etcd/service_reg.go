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

type ServiceInfo struct {
	Ip string
}

type Service struct {
	Name string
	Info ServiceInfo
	Cli  *clientv3.Client
	Stop  chan error
	Leaseid   clientv3.LeaseID
}

func ServiceRegInitClient(name string, info ServiceInfo) (*Service, error){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: conf.AppConfig.ETCD.EndPoints,
		DialTimeout: time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	service := &Service{Name: name, Info: info, Cli: cli, Stop: make(chan error)}
	return service, nil
}

func (service *Service) Start() error{
	ch, err := service.KeepAlive()
	if err != nil {
		log.SugarLogger.Error("etcd error:", err.Error())
		return err
	}
	for {
		select {
		case err <- service.Stop:
			service.Revoke()
			return err
		case <- service.Cli.Ctx().Done():
			return errors.New("server close")
		case re, ok := <- ch:
			// todo
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
	key := "service/" + service.Name
	value, _ := json.Marshal(service.Info)
	resp, err := service.Cli.Grant(context.TODO(), 5)
	if err != nil {
		return nil, err
	}
	_, err = service.Cli.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}
	service.Leaseid = resp.ID
	return service.Cli.KeepAlive(context.TODO(), resp.ID)
}


func (service *Service)Revoke() error{
	_, err := service.Cli.Revoke(context.TODO(), service.Leaseid)
	if err != nil {
		log.SugarLogger.Error("err:", err.Error())
	}
	return nil
}

func (service *Service) StopService()  {
	service.Stop <- nil
}


