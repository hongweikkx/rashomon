package util

import (
	"fmt"
	"testing"
)

func TestOrderMap(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "a"
	m["b"] = "b"
	m["c"] = "c"
	m["d"] = "d"
	m["e"] = "e"
	orderM := NewOrderMap(func(a, b interface{}) bool {
		return a.(string) < b.(string)
	})
	for k, v := range m {
		orderM.Add(k, v, k)
	}
	orderM.Delete("b")
	for _, v := range orderM.Iter() {
		fmt.Println(v.(KV).K)
	}
}
