package types

import (
	"context"
	"luna-backend/common"
	"luna-backend/db/internal/tables"
	"luna-backend/errors"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type MigrationQueries struct {
	Tx           pgx.Tx
	Context      context.Context
	Logger       *logrus.Entry
	CommonConfig *common.CommonConfig
	Tables       *tables.Tables
	Runner       func(*MigrationQueries, *common.Version) *errors.ErrorTrace
}

func (q *MigrationQueries) RunMigrations(lastVersion *common.Version) *errors.ErrorTrace {
	return q.Runner(q, lastVersion)
}
