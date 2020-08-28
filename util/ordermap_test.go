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
		return a.(int) < b.(int)
	})
	for k, v := range m {
		orderM.Add(k, v, MD5INT32(k))
	}
	for _, v := range orderM.Iter() {
		h := v.(KV).K.(string)
		fmt.Println("ret:", h, " ", MD5INT32(h))
	}
	orderM.Delete("b")
	for _, v := range orderM.Iter() {
		fmt.Println(v.(KV).K)
	}
}
