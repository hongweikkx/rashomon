package load_balance

import "errors"

type RoundRobinAL struct {
	Current int
	Keys    []string
}

func (al *RoundRobinAL) Init() {
	al.Current = -1
}

func (al *RoundRobinAL) ADD(server Server) {
}

func (al *RoundRobinAL) DELETE(key string) {
}

func (al *RoundRobinAL) GetNext(str string) (string, error) {
	if len(al.Keys) == 0 {
		return "", errors.New("no valid server to use")
	}
	cur := al.Current
	next := (cur + 1) % len(al.Keys)
	al.Current = next
	return al.Keys[al.Current], nil
}
