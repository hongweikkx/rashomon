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

// todo
func (al *ConsistentHashAL)GetNext() int {
	//for _, v := range al.Nodes.Iter() {
	//}
	return -1
}
