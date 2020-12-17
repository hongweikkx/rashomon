package load_balance

import (
	"testing"
)

func TestWeightedRoundRobinAL(t *testing.T) {
	nodes1 := WNode{
		Key:           "a",
		Weight:        5,
		CurrentWeight: 0,
	}
	nodes2 := WNode{
		Key:           "b",
		Weight:        1,
		CurrentWeight: 0,
	}
	nodes3 := WNode{
		Key:           "c",
		Weight:        1,
		CurrentWeight: 0,
	}
	nodes := []WNode{nodes2, nodes1, nodes3}
	w := WeightedRoundRobinAL{Nodes: nodes}
	res := []string{"a", "a", "b", "a", "c", "a", "a", "a", "a", "b"}
	for i := 0; i < 10; i++ {
		a, err := w.GetNext("")
		if err != nil {
			t.Error("TestWeightedRoundRobinAL:", err.Error())
		} else if a != res[i] {
			t.Error("TestWeightedRoundRobinAL", a)
		}
	}
}
