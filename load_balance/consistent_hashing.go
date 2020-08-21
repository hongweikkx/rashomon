package load_balance

import "net/url"

type CHNode struct {
	Url *url.URL
}

type ConsistentHashAL struct {
	Nodes []CHNode
}
func (al *ConsistentHashAL)GetNext() int {
	return -1
}
