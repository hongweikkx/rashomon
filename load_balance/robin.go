package load_balance

func (serverPool *ServerPool)GetNextWithRR() int{
	cur := serverPool.Current
	return (cur + 1)%len(serverPool.Servers)
}
