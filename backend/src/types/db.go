package types

import (
	"context"
	"io"
	"luna-backend/errors"
	"time"
)

type EventDatabaseEntry struct {
	Id       ID     `db:"id" encrypted:"false"`
	Calendar ID     `db:"calendar" encrypted:"false"`
	Settings []byte `db:"settings" encrypted:"false"`
}

type EventExtendedDatabaseEntry struct {
	Id          ID     `db:"id" encrypted:"false"`
	Calendar    ID     `db:"calendar" encrypted:"false"`
	Settings    []byte `db:"settings" encrypted:"false"`
	Title       string `db:"title" encrypted:"false"`
	Description string `db:"description" encrypted:"false"`
	Color       []byte `db:"color" encrypted:"false"`
	Overridden  bool   `db:"overridden" encrypted:"false"`
}

type CalendarDatabaseEntry struct {
	Id       ID     `db:"id" encrypted:"false"`
	Source   ID     `db:"source" encrypted:"false"`
	Settings []byte `db:"settings" encrypted:"false"`
}

type CalendarExtendedDatabaseEntry struct {
	Id          ID     `db:"id" encrypted:"false"`
	Source      ID     `db:"source" encrypted:"false"`
	Settings    []byte `db:"settings" encrypted:"false"`
	Title       string `db:"title" encrypted:"false"`
	Description string `db:"description" encrypted:"false"`
	Color       []byte `db:"color" encrypted:"false"`
	Overridden  bool   `db:"overridden" encrypted:"false"`
}

type SourceDatabaseEntry struct {
	Id       ID     `db:"id" encrypted:"false"`
	UserId   ID     `db:"userid" encrypted:"false"`
	Name     string `db:"name" encrypted:"false"`
	Type     string `db:"type" encrypted:"false"`
	Settings []byte `db:"settings" encrypted:"false"`
	AuthType string `db:"auth_type" encrypted:"true"`
	Auth     []byte `db:"auth" encrypted:"true"`
}

// Subset of database queries required for protocol implementations
// Required to avoid circular dependencies
type DatabaseQueries interface {
	GetContext() context.Context

	GetFilecache(file File) (string, io.Reader, *time.Time, *errors.ErrorTrace)
	SetFilecache(file File, content io.Reader) *errors.ErrorTrace
	SetFilecacheWithoutId(file File, content io.Reader) (ID, *errors.ErrorTrace)
	DeleteFilecache(file File) *errors.ErrorTrace

	SetCalendarOverrides(calendarId ID, name string, desc string, color *Color) *errors.ErrorTrace
	DeleteCalendarOverrides(calendarId ID) *errors.ErrorTrace
	SetEventOverrides(eventId ID, name string, desc string, color *Color) *errors.ErrorTrace
	DeleteEventOverrides(eventId ID) *errors.ErrorTrace
}
