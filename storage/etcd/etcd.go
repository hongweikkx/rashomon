package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"go.etcd.io/etcd/clientv3"
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

func New() (*Master, error){
	watchKey := conf.AppConfig.Storage.ETCD.WatchPrix
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.AppConfig.Storage.ETCD.EndPoints,
		DialTimeout: time.Duration(conf.AppConfig.Storage.ETCD.DailTimeout) * time.Second,
		Username: conf.AppConfig.Storage.ETCD.User,
		Password: conf.AppConfig.Storage.ETCD.Password,
	})
	if err != nil {
		return nil, err
	}
	// init
	master := &Master { Cli: cli}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := cli.Get(ctx, watchKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range res.Kvs {
		master.Add(v.Key, v.Value)
	}
	// watch
	go master.WatchWorkers(watchKey)
	return master, nil
}


func (master *Master)WatchWorkers(key string) {
	clientv3.NewWatcher(master.Cli)
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


func (this *Master)Stop() {
	err := this.Cli.Close()
	if err != nil {
		log.SugarLogger.Error("store cli close error:", err.Error())
	}
}
