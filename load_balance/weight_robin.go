package load_balance

import (
	"errors"
	"math"
)

type WNode struct {
	Key           string
	Weight        int
	CurrentWeight int
}

type WeightedRoundRobinAL struct {
	Nodes []WNode
}

func (al *WeightedRoundRobinAL) Init() {
}

func (al *WeightedRoundRobinAL) ADD(server Server) {
	al.Nodes = append(al.Nodes, WNode{Key: server.Key, Weight: server.Weight, CurrentWeight: 0})
}

func (al *WeightedRoundRobinAL) DELETE(key string) {
	for k, v := range al.Nodes {
		if v.Key == key {
			l := len(al.Nodes)
			if k == l-1 {
				al.Nodes = al.Nodes[:l-1]
			} else {
				al.Nodes = append(al.Nodes[:k], al.Nodes[k+1:]...)
			}
		}
	}
}

func (al *WeightedRoundRobinAL) GetNext(_str string) (string, error) {
	if len(al.Nodes) <= 0 {
		return "", errors.New("no valid server to use")
	}
	// 加weight
	// 选择
	maxWeight := math.MinInt32
	index := -1
	allWeight := 0
	for k, v := range al.Nodes {
		v.CurrentWeight += v.Weight
		al.Nodes[k] = v
		if v.CurrentWeight > maxWeight {
			maxWeight = v.CurrentWeight
			index = k
		}
		allWeight += v.Weight
	}
	// 替代
	if index != -1 {
		al.Nodes[index].CurrentWeight -= allWeight
		//fmt.Println(allWeight)
		//for _, v := range al.Nodes {
		//	fmt.Println(v.Key,v.CurrentWeight)
		//}
		key := al.Nodes[index].Key
		return key, nil
	}
	return "", errors.New("no valid server to use")
}
