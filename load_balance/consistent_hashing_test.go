package load_balance

import (
	"github.com/hongweikkx/rashomon/util"
	"testing"
)

func TestConsistentHashAL(t *testing.T) {
	keys := []string{"a", "b", "c", "d", "e"}
	orderM := util.NewOrderMap(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	})
	for _, key := range keys {
		md5 := util.MD5INT32(key)
		orderM.Add(key, CHNode{
			Key: key,
			Num: md5,
		}, md5)
	}
	w := &ConsistentHashAL{Nodes: orderM}
	res := []string{"d","c", "a", "e", "b"}
	for i := 0; i< len(keys); i++ {
		a, err := w.GetNext(keys[i])
		if err != nil {
			t.Error("TestWeightedRoundRobinAL:", err.Error())
		} else if a != res[i] {
			t.Error("TestWeightedRoundRobinAL", a)
		}
	}
}