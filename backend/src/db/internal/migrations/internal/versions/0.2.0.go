package versions

import (
	"fmt"
	"luna-backend/common"
	"luna-backend/db/internal/migrations/internal/registry"
	"luna-backend/db/internal/migrations/types"
)

func init() {
	registry.RegisterMigration(common.Ver(0, 2, 0), func(q *types.MigrationQueries) error {
		err := q.Tables.InitializeCalendarsTable()
		if err != nil {
			return fmt.Errorf("could not initialize calendars table: %v", err)
		}

		return nil
	})
}
