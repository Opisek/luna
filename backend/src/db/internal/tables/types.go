package tables

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Tables struct {
	Tx      pgx.Tx
	Context context.Context
	//logger       *logrus.Entry
	//commonConfig *common.CommonConfig
}
