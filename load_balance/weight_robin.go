package load_balance

import "net/url"

type node struct {
	Url *url.URL
	Weight int
	CurrentWeight int
}

type WeightedRoundRobinAL struct {
	Nodes []node
}

func (al WeightedRoundRobinAL)GetNext() int{
	if len(al.Nodes) <= 0 {
		return -1
	}
	// 加weight
	// 选择
	maxWeight := 0
	index := -1
	allWeight := 0
	for k, v := range al.Nodes {
		v.CurrentWeight += v.Weight
		if v.CurrentWeight > maxWeight {
			maxWeight = v.CurrentWeight
			index = k
			allWeight += v.Weight
		}
	}
	// 替代
	if index != -1 {
		al.Nodes[index].Weight -= allWeight
		urlIndex := al.Nodes[index].Url
		for k, v := range ServerPoolLB.Servers {
			if v.URL == urlIndex {
				return k
			}
		}
	}
	return -1
}
