package load_balance


type RoundRobinAL struct {
}

func (al RoundRobinAL)GetNext() int{
	cur := ServerPoolLB.Current
	return (cur + 1)%len(ServerPoolLB.Servers)
}
