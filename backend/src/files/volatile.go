package files

import (
	"bytes"
	"io"
	"luna-backend/errors"
	"luna-backend/types"
)

type VolatileFile struct {
	name    string
	content []byte
}

func NewVolatileFile(name string, content []byte) *VolatileFile {
	return &VolatileFile{name: name, content: content}
}

func (file *VolatileFile) GetId() types.ID {
	panic("Volatile files do not have an ID")
}

func (file *VolatileFile) SetId(_ types.ID) {
	panic("Volatile files do not have an ID")
}

func (file *VolatileFile) GetName(_ types.DatabaseQueries) string {
	return file.name
}

func (file *VolatileFile) GetContent(_ types.DatabaseQueries) (io.Reader, *errors.ErrorTrace) {
	return bytes.NewReader(file.content), nil
}

func (file *VolatileFile) GetBytes(_ types.DatabaseQueries) ([]byte, *errors.ErrorTrace) {
	return file.content, nil
}
