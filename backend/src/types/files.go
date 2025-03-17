package types

import (
	"io"
	"luna-backend/errors"
)

type File interface {
	GetId() ID
	SetId(id ID)
	GetContent(q DatabaseQueries) (io.Reader, *errors.ErrorTrace)
}
