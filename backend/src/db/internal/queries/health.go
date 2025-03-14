package queries

import "golang.org/x/net/context"

func (q *Queries) CheckHealth() error {
	return q.Tx.Conn().Ping(context.TODO())
}
