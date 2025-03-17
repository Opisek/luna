package migrations

import (
	"luna-backend/common"
	"luna-backend/db/internal/migrations/internal/registry"
	_ "luna-backend/db/internal/migrations/internal/versions"
	"luna-backend/db/internal/migrations/types"
	"luna-backend/db/internal/tables"
	"luna-backend/errors"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func NewMigrationQueries(tx pgx.Tx, logger *logrus.Entry, commonConfig *common.CommonConfig, tables *tables.Tables) *types.MigrationQueries {
	return &types.MigrationQueries{
		Tx:           tx,
		Logger:       logger,
		CommonConfig: commonConfig,
		Tables:       tables,
		Runner:       runMigrations,
	}
}

func runMigrations(q *types.MigrationQueries, lastVersion *common.Version) *errors.ErrorTrace {

	migrations := registry.GetMigrations(*lastVersion)

	for _, migration := range migrations {
		q.Logger.Infof("running migration %s", migration.Ver.String())
		err := migration.Fun(q)
		if err != nil {
			return err.
				Append(errors.LvlDebug, "Error running migration for %s", migration.Ver.String())
		}
	}

	q.Logger.Infof("migrations up to date")
	return nil
}
