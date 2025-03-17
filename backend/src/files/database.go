package files

import (
	"bytes"
	"io"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
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

func NewDatabaseFileFromContent(content io.Reader, q types.DatabaseQueries) (*DatabaseFile, *errors.ErrorTrace) {
	buf, err := io.ReadAll(content)
	file := &DatabaseFile{content: buf}
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not read file content").
			Append(errors.LvlPlain, "Could not upload file")
	}
	_, tr := q.SetFilecacheWithoutId(file, bytes.NewReader(buf))
	if tr != nil {
		return nil, tr.Append(errors.LvlPlain, "Could not upload file")
	}
	return file, nil
}

func (file *DatabaseFile) GetId() types.ID {
	return file.id
}

func (file *DatabaseFile) SetId(id types.ID) {
	file.id = id
}

func (file *DatabaseFile) fetchContentFromDatabase(q types.DatabaseQueries) (io.Reader, *time.Time, *errors.ErrorTrace) {
	content, date, err := q.GetFilecache(file)
	if err != nil {
		return nil, nil, err.
			Append(errors.LvlDebug, "Could not get file %v from the database", file.GetId()).
			AltStr(errors.LvlPlain, "Could not get file from the database")
	}
	return content, date, nil
}

func (file *DatabaseFile) GetContent(q types.DatabaseQueries) (io.Reader, *errors.ErrorTrace) {
	if file.content != nil {
		return bytes.NewReader(file.content), nil
	}

	reader, _, err := file.fetchContentFromDatabase(q)
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not read contents of file %v", file.GetId()).
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}

	return reader, nil
}
