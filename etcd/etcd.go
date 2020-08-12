package etcd

import (
	"context"
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"go.etcd.io/etcd/client"
	"sync"
	"time"
)


type Master struct{
	Members sync.Map
	KeysAPI client.KeysAPI
}

type Member struct {
	InGroup bool
	Value string
}


func New() error{
	cli, err := client.New(client.Config{
		Endpoints:   conf.AppConfig.ETCD.EndPoints,
		Transport: client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second,
	})
	if err != nil {
		return err
	}
	master := &Master {
		KeysAPI: client.NewKeysAPI(cli),
	}
	go master.WatchWorkers()
	return nil
}

func (master *Master)WatchWorkers() {
	watcher := master.KeysAPI.Watcher("workers/", &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.SugarLogger.Error("watch error:", err.Error())
			break
		}
		// action: get, set, delete, update, create, compareAndSwap, compareAndDelete and expire.
		// todo 我不认为这样做是对的  因为首先expire的使用 在blog上就不是很好
		// todo 其次这次 action没有都包含进去， 比如compareAnd 这些是？
		// todo 看下k8s是怎么做的吧
		switch res.Action {
		case "expire":
			err := master.Update(res.PrevNode.Key, false)
			if err != nil {
				log.SugarLogger.Error("service member error:", err.Error())
			}
		case "delete":
			master.Delete(res.Node.Key)
		case "set", "create":
			master.Add(res.Node.Key, &Member{InGroup: true, Value: res.Node.Value})
		}
	}
}


func (this *Master)Update(key string, inGroup bool) error{
	if v, ok := this.Members.Load(key); ok {
		member := v.(Member)
		member.InGroup = inGroup
		this.Members.Store(key, member)
		return nil
	}else {
		return errors.New("member:" + key + "is not exist")
	}
}

func (this *Master)Add(key string, member *Member) {
	member.InGroup = true
	this.Members.Store(key, member)
}

func (this *Master)Delete(key string) {
	this.Delete(key)
}