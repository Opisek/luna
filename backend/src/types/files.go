package types

import (
	"io"
	"time"
)

// We cannot import db in here so instead of using the whole *db.Queries, we create a subset interface
type FileQueries interface {
	GetFilecache(file File) (io.Reader, *time.Time, error)
	SetFilecache(file File, content io.Reader) error
	DeleteFilecache(file File) error
}

type File interface {
	GetId() ID
	GetContent(q FileQueries) (io.Reader, error)
}
