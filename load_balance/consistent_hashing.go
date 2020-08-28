package load_balance

import (
	"github.com/hongweikkx/rashomon/util"
)

type CHNode struct {
	Key string
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
	md5 := util.MD5INT32(server.Key)
	al.Nodes.Add(server.Key, CHNode{
		Key: server.Key,
		Num:     md5,
	}, md5)
}

func (al *ConsistentHashAL) DELETE(address string) {
	al.Nodes.Delete(address)
}

func (al *ConsistentHashAL)GetNext(str string) (string, error) {
	firstKV, err := al.Nodes.First()
	if err != nil {
		return "", err
	}
	key := firstKV.(util.KV).K.(string)
	md5 := util.MD5INT32(str)
	for _, node := range al.Nodes.Iter() {
		chNode := node.(util.KV).V.(CHNode)
		if md5 < chNode.Num {
			key = chNode.Key
			break
		}
	}

	return key, nil
}
