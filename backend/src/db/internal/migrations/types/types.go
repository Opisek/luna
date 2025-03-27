package types

import (
	"context"
	"luna-backend/config"
	"luna-backend/db/internal/tables"
	"luna-backend/errors"
	"luna-backend/types"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type MigrationQueries struct {
	Tx           pgx.Tx
	Context      context.Context
	Logger       *logrus.Entry
	CommonConfig *config.CommonConfig
	Tables       *tables.Tables
	Runner       func(*MigrationQueries, *types.Version) *errors.ErrorTrace
}

func (q *MigrationQueries) RunMigrations(lastVersion *types.Version) *errors.ErrorTrace {
	return q.Runner(q, lastVersion)
}
