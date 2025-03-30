package migrations

import (
	"context"
	"luna-backend/config"
	"luna-backend/db/internal/migrations/internal/registry"
	_ "luna-backend/db/internal/migrations/internal/versions"
	migrationTypes "luna-backend/db/internal/migrations/types"
	"luna-backend/db/internal/tables"
	"luna-backend/errors"
	"luna-backend/types"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func NewMigrationQueries(tx pgx.Tx, context context.Context, logger *logrus.Entry, commonConfig *config.CommonConfig, tables *tables.Tables) *migrationTypes.MigrationQueries {
	return &migrationTypes.MigrationQueries{
		Tx:           tx,
		Context:      context,
		Logger:       logger,
		CommonConfig: commonConfig,
		Tables:       tables,
		Runner:       runMigrations,
	}
}

func runMigrations(q *migrationTypes.MigrationQueries, lastVersion *types.Version) *errors.ErrorTrace {

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
