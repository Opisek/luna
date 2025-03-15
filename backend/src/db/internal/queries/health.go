package queries

func (q *Queries) CheckHealth() error {
	return q.Tx.Conn().Ping(q.Context)
}
