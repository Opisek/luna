package parsing

import (
	"luna-backend/errors"
	"luna-backend/interface/primitives"
	"luna-backend/types"
)

// Little hack so we can use methods in interface/parsing.
// Those methods cannot be moved here or referenced directly due to a circular dependency.
type PrimitivesParser interface {
	ParseSource(entry *types.SourceDatabaseEntry) (primitives.Source, *errors.ErrorTrace)
	ParseCalendarSettings(sourceType string, settings []byte) (primitives.CalendarSettings, *errors.ErrorTrace)
	ParseEventSettings(sourceType string, settings []byte) (primitives.EventSettings, *errors.ErrorTrace)
}
