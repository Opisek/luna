package tables

import (
	"github.com/jackc/pgx/v5"
)

type Tables struct {
	Tx pgx.Tx
	//logger       *logrus.Entry
	//commonConfig *common.CommonConfig
}
