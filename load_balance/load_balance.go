package load_balance

import (
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"net/url"
	"sync"
)

const WeightedRoundRobin = 1
const RoundRobin = 2
const ConsistentHashing = 3

type Server struct {
	URL   *url.URL
}

type ServerPool struct {
	Servers []Server
	lock  sync.RWMutex
	LoadBalance LoadBalanceAPI
}

type LoadBalanceAPI interface {
	// must hold the lock
	GetNext() int
}

var ServerPoolLB ServerPool

// todo 1. simple factory
func init() {
	ServerPoolLB =
	ServerPool{
		Servers:     nil,
		lock:        sync.RWMutex{},
		LoadBalance: NewLBAPI(),
	}
}

func NewLBAPI() LoadBalanceAPI{
	switch conf.AppConfig.LoadBalance.Algorithm  {
	case WeightedRoundRobin:
		return &WeightedRoundRobinAL{}
	case ConsistentHashing:
		return &ConsistentHashAL{}
	default:
		return &RoundRobinAL{}
	}
}

func (serverPool *ServerPool)UpdateServer(s Server) {
	defer serverPool.lock.Unlock()
	serverPool.lock.Lock()
	if index := serverPool.isConlict(s.URL); index != -1 {
		serverPool.Servers[index] = s
		return
	}else {
		serverPool.Servers = append(serverPool.Servers, s)
	}
}

func (serverPool *ServerPool)GetNext() (*Server, error){
	defer serverPool.lock.Unlock()
	serverPool.lock.Lock()
	index := serverPool.LoadBalance.GetNext()
	if index == -1 {
		log.SugarLogger.Error("none server can use")
		return nil, errors.New("none server can use")
	}
	return &serverPool.Servers[index], nil
}


// must hold the lock
func (serverPool *ServerPool)isConlict(url *url.URL) int{
	for k, v := range serverPool.Servers {
		if v.URL == url {
			return k
		}
	}
	return -1
}


