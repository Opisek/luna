package parsing

import (
	"luna-backend/errors"
	"luna-backend/types"
)

// Little hack so we can use methods in interface/parsing.
// Those methods cannot be moved here or referenced directly due to a circular dependency.
type PrimitivesParser interface {
	ParseSource(entry *types.SourceDatabaseEntry) (types.Source, *errors.ErrorTrace)
	ParseCalendarSettings(sourceType string, settings []byte) (types.CalendarSettings, *errors.ErrorTrace)
	ParseEventSettings(sourceType string, settings []byte) (types.EventSettings, *errors.ErrorTrace)
}
