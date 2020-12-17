package util

import (
	"crypto/md5"
	"math"
	"math/big"
)

func IntArrDelete(arr *[]int, index int) {
	l := len(*arr)
	if index >= l {
		return
	} else if index == l-1 {
		*arr = (*arr)[:l-1]
	} else {
		*arr = append((*arr)[:index], (*arr)[index+1:]...)
	}
}

func MD5INT32(str string) int {
	biMd5 := big.NewInt(0)
	h := md5.New()
	h.Write([]byte(str))
	biMd5.SetBytes(h.Sum(nil))

	biINT32MAX := big.NewInt(math.MaxInt32)

	ret := big.NewInt(0)
	ret.Mod(biMd5, biINT32MAX)
	return int(ret.Int64())
}
