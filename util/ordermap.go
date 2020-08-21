package util

import "sort"

type OrderKey struct {
	key     interface{}
	rankKey interface{}
}

type MapOrderKeys struct {
	sortedList []OrderKey
	lessFunc   func(a, b interface{}) bool
}

type OrderMap struct {
	sortedKeys MapOrderKeys
	mapI       map[interface{}]interface{}
}

type KV struct {
	K interface{}
	V interface{}
}

func NewOrderMap(f func(a, b interface{})bool) *OrderMap{
	mapOrderKeys := MapOrderKeys{lessFunc: f}
	orderMap := &OrderMap{sortedKeys: mapOrderKeys, mapI: make(map[interface{}]interface{})}
	return orderMap
}

func (om *OrderMap) Add(key interface{}, value interface{}, rankKey interface{}){
	if _, ok := om.mapI[key]; ok {
		om.mapI[key] = value
		for k, v := range om.sortedKeys.sortedList {
			if v.key == key {
				om.sortedKeys.sortedList[k].rankKey = rankKey
				break
			}
		}
		sort.Sort(om.sortedKeys)
	}else {
		om.mapI[key] = value
		om.sortedKeys.sortedList = append(om.sortedKeys.sortedList, OrderKey{
			key:     key,
			rankKey: rankKey,
		})
	}
}

func (om *OrderMap) Delete(key interface{}) {
	if _, ok := om.mapI[key]; ok {
		delete(om.mapI, key)
		for k, v := range om.sortedKeys.sortedList {
			if v.key == key {
				l := len(om.sortedKeys.sortedList)
				if k == l -1 {
					om.sortedKeys.sortedList = om.sortedKeys.sortedList[:l-1]
				}else {
					om.sortedKeys.sortedList = append(om.sortedKeys.sortedList[:k], om.sortedKeys.sortedList[k+1:]...)
				}
				break
			}
		}
	}
}

func (om *OrderMap) Iter() []interface{}{
	ret := []interface{}{}
	for _, v := range om.sortedKeys.sortedList {
		ret = append(ret, KV{v.key,om.mapI[v.key]})
	}
	return ret
}

func (keys MapOrderKeys) Len() int{
	return len(keys.sortedList)
}

func (keys MapOrderKeys) Swap(i, j int) {
	keys.sortedList[i], keys.sortedList[j] = keys.sortedList[j], keys.sortedList[i]
}

func (keys MapOrderKeys) Less(i, j int) bool{
	return keys.lessFunc(keys.sortedList[i], keys.sortedList[j])
}

