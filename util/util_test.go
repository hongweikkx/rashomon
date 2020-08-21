package util

import (
	"fmt"
	"testing"
)

func TestMD5INT(t *testing.T) {
	a := MD5INT32("hello, world")
	fmt.Print(a)
}

