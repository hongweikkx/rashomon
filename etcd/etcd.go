package etcd

import (
	"context"
	"encoding/json"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"go.etcd.io/etcd/clientv3"
	mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
	"time"
)


// todo 看看k8s的做法
type Master struct{
	Cli *clientv3.Client
	Members sync.Map
}

type ServiceInfo struct {
	Ip string
	Port string
}

type Member struct {
	EndPoint ServiceInfo
}

// 新建一个监控
func New() (*Master, error){
	watchKey := conf.AppConfig.ETCD.WatchPrix
	// etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.AppConfig.ETCD.EndPoints,
		DialTimeout: time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	// init
	res, err := cli.Get(context.Background(), watchKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	master := &Master { Cli: cli}
	for _, v := range res.Kvs {
		master.Add(v.Key, v.Value)
	}
	// watch
	go master.WatchWorkers(watchKey)
	return master, nil
}


func (master *Master)WatchWorkers(key string) {
	watchCh := master.Cli.Watch(context.Background(), key, clientv3.WithPrefix())
	for watchMsg := range watchCh {
		for _, event := range watchMsg.Events {
			switch event.Type {
			case mvccpb.PUT:
				master.Add(event.Kv.Key, event.Kv.Value)
			case mvccpb.DELETE:
				master.Delete(event.Kv.Key)
			}
		}

	}
}

func (this *Master)Add(keyByte []byte, serviceByte []byte) {
	key := string(keyByte)
	var service ServiceInfo
	err := json.Unmarshal(serviceByte, &service)
	if err != nil {
		log.SugarLogger.Error("Add error:" + key + err.Error())
		return
	}
	this.Members.Store(key, service)
}

func (this *Master)Delete(keyByte []byte) {
	this.Members.Delete(string(keyByte))
}

