package load_balance

import (
	"github.com/hongweikkx/rashomon/log"
	"sync"
)

const WeightedRoundRobin = 1
const RoundRobin = 2
const ConsistentHashing = 3

type Cluster struct {
	Servers []*Server
	LBStrategy int
	LoadBalance LoadBalanceAPI
	lock  sync.RWMutex
}

type Server struct {
	Key       string
	Address   string
	Weight    int
	ConsistentHashStr string
}

type LoadBalanceAPI interface {
	Init()
	// must hold the lock
	GetNext(string) (string, error)
	ADD(server Server)
	DELETE(address string)
}

func NewLBAPI(cluster *Cluster) LoadBalanceAPI{
	var api LoadBalanceAPI
	switch cluster.LBStrategy {
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

func (cluster *Cluster)UpdateServer(s *Server) {
	cluster.lock.Lock()
	defer cluster.lock.Unlock()
	if index := cluster.isExist(s.Key); index != -1 {
		cluster.Servers[index] = s
	}else {
		cluster.Servers = append(cluster.Servers, s)
	}
	cluster.LoadBalance.ADD(*s)
}

func (cluster *Cluster)DeleteServer(s Server) {
	cluster.lock.Lock()
	defer cluster.lock.Unlock()
	if index := cluster.isExist(s.Key); index != -1 {
		// delete
		l := len(cluster.Servers)
		if index == l -1 {
			cluster.Servers = cluster.Servers[:l-1]
		}else {
			cluster.Servers = append(cluster.Servers[:index], cluster.Servers[index+1:]...)
		}
		cluster.LoadBalance.DELETE(s.Key)
	}
}

func (cluster *Cluster)GetNext(key string) (string, error){
	defer cluster.lock.Unlock()
	cluster.lock.Lock()
	key, err := cluster.LoadBalance.GetNext(key)
	if err != nil {
		log.SugarLogger.Error("none server can use")
		return key, err
	}
	return key, nil
}


// must hold the lock
func (cluster *Cluster) isExist(key string) int{
	for k, v := range cluster.Servers {
		if v.Key == key {
			return k
		}
	}
	return -1
}



