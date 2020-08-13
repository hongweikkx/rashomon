package etcd

import (
	"context"
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/client"
	"sync"
	"time"
)


// todo  还是应该采用最新的etcd v3的实现来做。
// todo  先自己摸索下吧
// todo 先用网上的说法吧。 然后在看看k8s的做法

type Master struct{
	Cli *clientv3.Client
	Members []Member
	lock  sync.RWMutex
	WatchKey string
}

type Member struct {
	InGroup bool
	EndPoint EndPoint
}

type EndPoint struct {
	IP string
	Port int
}

// 新建一个监控
func New(WatchKey string, endPoints []EndPoint) (*Master, error){
	// etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.AppConfig.ETCD.EndPoints,
		DialTimeout: time.Duration(conf.AppConfig.ETCD.DailTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	// init members
	members := []Member{}
	for _, endPoint := range endPoints {
		members = append(members, Member{InGroup: false, EndPoint: endPoint})
	}
	// fill master
	master := &Master {
		Cli: cli,
		Members: members,
		WatchKey: WatchKey,
	}
	// watch
	go master.WatchWorkers(master.WatchKey)
	return master, nil
}


func (master *Master)WatchWorkers(key string) {
	watcher := master.KeysAPI.Watcher(key, &client.WatcherOptions{
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