package load_balance

import "testing"

func TestRoundRobinAL(t *testing.T) {
	keys := []string{"a", "a", "b", "a", "c", "a"}
	w := &RoundRobinAL{
		Current: 0,
		Keys:    keys,
	}
	res := []string{"a", "b", "a", "c", "a", "a", "a", "b"}
	for i := 0; i < len(res); i++ {
		a, err := w.GetNext("")
		if err != nil {
			t.Error("TestWeightedRoundRobinAL:", err.Error())
		} else if a != res[i] {
			t.Error("TestWeightedRoundRobinAL", a)
		}
	}
}
