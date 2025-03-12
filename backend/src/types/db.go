package types

import (
	"io"
	"time"
)

type EventDatabaseEntry struct {
	Id       ID     `db:"id" encrypted:"false"`
	Calendar ID     `db:"calendar" encrypted:"false"`
	Color    []byte `db:"color" encrypted:"false"`
	Settings []byte `db:"settings" encrypted:"false"`
}

type CalendarDatabaseEntry struct {
	Id       ID     `db:"id" encrypted:"false"`
	Source   ID     `db:"source" encrypted:"false"`
	Color    []byte `db:"color" encrypted:"false"`
	Settings []byte `db:"settings" encrypted:"false"`
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

// Subset of database queries required for e.g. file caching
// Required to avoid circular dependencies
type DatabaseQueries interface {
	GetFilecache(file File) (io.Reader, *time.Time, error)
	SetFilecache(file File, content io.Reader) error
	SetFilecacheWithoutId(file File, content io.Reader) (ID, error)
	DeleteFilecache(file File) error
}
