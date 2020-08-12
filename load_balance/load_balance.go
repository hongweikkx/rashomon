package load_balance

import (
	"errors"
	"github.com/hongweikkx/rashomon/conf"
	"github.com/hongweikkx/rashomon/log"
	"net/url"
	"sync"
)

type Server struct {
	URL   *url.URL
	Weight int
	Alive  bool
}

type ServerPool struct {
	Servers []Server
	Current int
	lock  sync.RWMutex
}


const WeightedRoundRobin = 1
const RoundRobin = 2
const ConsistentHashing = 3

func (serverPool *ServerPool)ReplaceServer(s Server) {
	defer serverPool.lock.Unlock()
	serverPool.lock.Lock()
	if index := isConlict(serverPool.Servers, s.URL); index != -1 {
		serverPool.Servers[index] = s
		return
	}else {
		serverPool.Servers = append(serverPool.Servers, s)
	}
}

func (serverPool *ServerPool)GetNext() (Server, error){
	defer serverPool.lock.RUnlock()
	serverPool.lock.RLock()
	index := serverPool.Current
	switch conf.AppConfig.LoadBalance.Algorithm {
	case WeightedRoundRobin:
		index = serverPool.GetNextWithWRR()
	case ConsistentHashing:
		index =  serverPool.GetNextWithCH()
	default:
		index =  serverPool.GetNextWithRR()
	}
	if index == -1 {
		log.SugarLogger.Error("none server can use")
		return serverPool.Servers[serverPool.Current], errors.New("none server can use")
	}
	serverPool.Current = index
	return serverPool.Servers[index], nil
}


func isConlict(servers []Server, url *url.URL) int{
	for i:=0; i< len(servers); i++ {
		if servers[i].URL == url {
			return i
		}
	}
	return -1
}


