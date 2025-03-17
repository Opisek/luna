package types

import (
	"io"
	"luna-backend/errors"
)

type File interface {
	GetId() ID
	SetId(id ID)
	GetName(q DatabaseQueries) string
	GetContent(q DatabaseQueries) (io.Reader, *errors.ErrorTrace)
	GetBytes(q DatabaseQueries) ([]byte, *errors.ErrorTrace)
}
