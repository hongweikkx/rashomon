package load_balance

import (
	"github.com/hongweikkx/rashomon/util"
)

type CHNode struct {
	Address string
	Num int
}

// url ->
type ConsistentHashAL struct {
	Nodes *util.OrderMap
}

func (al *ConsistentHashAL)Init() {
	nodes := util.NewOrderMap(func(a, b interface{}) bool{
		return a.(int) < b.(int)
	})
	al.Nodes = nodes
}

func (al *ConsistentHashAL) ADD(server Server) {
	md5 := util.MD5INT32(server.Address)
	al.Nodes.Add(server.Address, CHNode{
		Address: server.Address,
		Num:     md5,
	}, md5)
}

func (al *ConsistentHashAL) DELETE(address string) {
	al.Nodes.Delete(address)
}

func (al *ConsistentHashAL)GetNext(str string) int {
	firstKV, err := al.Nodes.First()
	if err != nil {
		return -1
	}
	address := firstKV.(util.KV).K
	md5 := util.MD5INT32(str)
	for _, node := range al.Nodes.Iter() {
		chNode := node.(util.KV).V.(CHNode)
		if md5 < chNode.Num {
			address = chNode.Address
			break
		}
	}

	for k, v := range ServerPoolLB.Servers {
		if v.Address == address {
			return k
		}
	}
	return -1
}
