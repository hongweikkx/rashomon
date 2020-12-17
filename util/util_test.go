package util

import (
	"testing"
)

func TestMD5INT(t *testing.T) {
	a := MD5INT32("hello, world")
	if a != 128720900 {
		t.Errorf("err:%d", a)
	}
}

func TestIntArrDelete(t *testing.T) {
	arr := []int{1, 2, 3}
	IntArrDelete(&arr, 0)
	if arr[0] != 2 || arr[1] != 3 {
		t.Errorf("err:%d", arr)
	}
}
