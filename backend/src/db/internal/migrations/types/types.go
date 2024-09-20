package types

import (
	"luna-backend/common"
	"luna-backend/db/internal/tables"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type MigrationQueries struct {
	Tx           pgx.Tx
	Logger       *logrus.Entry
	CommonConfig *common.CommonConfig
	Tables       *tables.Tables
	Runner       func(*MigrationQueries, *common.Version) error
}

func (q *MigrationQueries) RunMigrations(lastVersion *common.Version) error {
	return q.Runner(q, lastVersion)
}
