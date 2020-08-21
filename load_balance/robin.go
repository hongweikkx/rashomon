package load_balance


type RoundRobinAL struct {
	Current int
}

func (al *RoundRobinAL)GetNext() int{
	cur := al.Current
	next := (cur + 1)%len(ServerPoolLB.Servers)
	al.Current = next
	return next
}
