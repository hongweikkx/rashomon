package load_balance

type WNode struct {
	Address string
	Weight int
	CurrentWeight int
}

type WeightedRoundRobinAL struct {
	Nodes []WNode
}

func (al *WeightedRoundRobinAL) Init() {
}

func (al *WeightedRoundRobinAL) ADD(server Server) {
	al.Nodes = append(al.Nodes, WNode{Address: server.Address, Weight: server.Weight, CurrentWeight: 0})
}

func (al *WeightedRoundRobinAL) DELETE(address string) {
	for k, v := range al.Nodes {
		if v.Address == address {
			l := len(al.Nodes)
			if k == l -1 {
				al.Nodes = al.Nodes[:l-1]
			}else {
				al.Nodes = append(al.Nodes[:k], al.Nodes[k+1:]...)
			}
		}
	}
}

func (al *WeightedRoundRobinAL)GetNext() int{
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
		address := al.Nodes[index].Address
		for k, v := range ServerPoolLB.Servers {
			if v.Address == address {
				return k
			}
		}
	}
	return -1
}
