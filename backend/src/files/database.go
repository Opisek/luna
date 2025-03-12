package files

import (
	"bytes"
	"fmt"
	"io"
	"luna-backend/types"
	"time"
)

// Implements types.File
type DatabaseFile struct {
	id      types.ID
	content []byte
}

func NewDatabaseFile(id types.ID) *DatabaseFile {
	return &DatabaseFile{id: id}
}

func NewDatabaseFileFromContent(content []byte, q types.FileQueries) (*DatabaseFile, error) {
	file := &DatabaseFile{content: content}
	_, err := q.SetFilecacheWithoutId(file, bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("could not create file: %w", err)
	}
	return file, nil
}

func (file *DatabaseFile) GetId() types.ID {
	return file.id
}

func (file *DatabaseFile) SetId(id types.ID) {
	file.id = id
}

func (file *DatabaseFile) fetchContentFromDatabase(q types.FileQueries) (io.Reader, *time.Time, error) {
	content, date, err := q.GetFilecache(file)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get file content from database: %w", err)
	}
	return content, date, nil
}

func (file *DatabaseFile) GetContent(q types.FileQueries) (io.Reader, error) {
	if file.content != nil {
		return bytes.NewReader(file.content), nil
	}

	reader, _, err := file.fetchContentFromDatabase(q)
	if err != nil {
		return nil, fmt.Errorf("could not get file content: %w", err)
	}

	return reader, nil
}
