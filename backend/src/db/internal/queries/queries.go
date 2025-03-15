package queries

import (
	"context"
	"luna-backend/common"
	"luna-backend/db/internal/parsing"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Queries struct {
	Tx               pgx.Tx
	Context          context.Context
	Logger           *logrus.Entry
	CommonConfig     *common.CommonConfig
	PrimitivesParser *parsing.PrimitivesParser
}

func (q *Queries) GetContext() context.Context {
	return q.Context
}
