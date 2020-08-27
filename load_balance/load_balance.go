package load_balance

import (
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"sync"
)

const WeightedRoundRobin = 1
const RoundRobin = 2
const ConsistentHashing = 3

type Server struct {
	Address   string
	Weight    int
}

type ServerPool struct {
	Servers []Server
	lock  sync.RWMutex
	LoadBalance LoadBalanceAPI
}

type LoadBalanceAPI interface {
	Init()
	// must hold the lock
	GetNext(string) int
	ADD(server Server)
	DELETE(address string)
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
	var api LoadBalanceAPI
	switch conf.AppConfig.LoadBalance.Algorithm  {
	case WeightedRoundRobin:
		api = &WeightedRoundRobinAL{}
	case ConsistentHashing:
		api = &ConsistentHashAL{}
	default:
		api = &RoundRobinAL{}
	}
	api.Init()
	return api
}

func (serverPool *ServerPool)UpdateServer(s Server) {
	serverPool.lock.Lock()
	defer serverPool.lock.Unlock()
	if index := serverPool.isExist(s.Address); index != -1 {
		serverPool.Servers[index] = s
	}else {
		serverPool.Servers = append(serverPool.Servers, s)
	}
	serverPool.LoadBalance.ADD(s)
}

func (ServerPool *ServerPool)DeleteServer(s Server) {
	ServerPool.lock.Lock()
	defer ServerPool.lock.Unlock()
	if index := ServerPool.isExist(s.Address); index != -1 {
		// delete
		l := len(ServerPool.Servers)
		if index == l -1 {
			ServerPool.Servers = ServerPool.Servers[:l-1]
		}else {
			ServerPool.Servers = append(ServerPool.Servers[:index], ServerPool.Servers[index+1:]...)
		}
		ServerPool.LoadBalance.DELETE(s.Address)
	}
}

func (serverPool *ServerPool)GetNext() (*Server, error){
	defer serverPool.lock.Unlock()
	serverPool.lock.Lock()
	// todo the args
	index := serverPool.LoadBalance.GetNext("")
	if index == -1 {
		log.SugarLogger.Error("none server can use")
		return nil, errors.New("none server can use")
	}
	return &serverPool.Servers[index], nil
}


// must hold the lock
func (serverPool *ServerPool) isExist(address string) int{
	for k, v := range serverPool.Servers {
		if v.Address == address {
			return k
		}
	}
	return -1
}


