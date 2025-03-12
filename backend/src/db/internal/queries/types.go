package queries

import (
	"luna-backend/common"
	"luna-backend/db/internal/parsing"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Queries struct {
	Tx               pgx.Tx
	Logger           *logrus.Entry
	CommonConfig     *common.CommonConfig
	PrimitivesParser *parsing.PrimitivesParser
}
