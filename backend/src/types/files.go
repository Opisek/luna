package types

import (
	"io"
)

type File interface {
	GetId() ID
	SetId(id ID)
	GetContent(q DatabaseQueries) (io.Reader, error)
}
