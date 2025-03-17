package queries

import (
	"luna-backend/errors"
	"net/http"
)

func (q *Queries) CheckHealth() *errors.ErrorTrace {
	switch q.Tx.Conn().Ping(q.Context) {
	case nil:
		return nil
	default:
		return errors.New().Status(http.StatusServiceUnavailable)
	}
}
