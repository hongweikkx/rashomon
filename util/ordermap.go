package util

type OrderKey struct {
	Key interface{}
	RankKey interface{}
}

type OrderMap struct {
	SortedList []OrderKey
	Map map[interface{}]interface{}
	RankFunc func(i, j int) bool
}

func New(f func(i, j int)bool) *OrderMap{
	orderMap := &OrderMap{RankFunc: f}
	return orderMap
}

func (om *OrderMap) Add(key interface{}, value interface{}, rankKey interface{}){
	if _, ok := om.Map[key]; ok {
		// exit 替换key
		for _, v := range
	}else {
		om.Map[key] = value
		om.SortedList = append(om.SortedList, OrderKey{
			Key:     key,
			RankKey: rankKey,
		})
	}
}

