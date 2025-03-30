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
	name    string
	content []byte
}

func GetDatabaseFile(id types.ID) *DatabaseFile {
	return &DatabaseFile{id: id}
}

func NewDatabaseFileFromContent(name string, content io.Reader, user types.ID, q types.DatabaseQueries) (*DatabaseFile, *errors.ErrorTrace) {
	buf, err := io.ReadAll(content)
	file := &DatabaseFile{name: name, content: buf}
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not read file content").
			Append(errors.LvlPlain, "Could not upload file")
	}
	_, tr := q.SetFilecacheWithoutId(file, bytes.NewReader(buf), user)
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

func (file *DatabaseFile) GetName(q types.DatabaseQueries) string {
	if file.content == nil {
		file.GetContent(q)
	}
	return file.name
}

func (file *DatabaseFile) fetchContentFromDatabase(q types.DatabaseQueries) (string, io.Reader, *time.Time, *errors.ErrorTrace) {
	name, content, date, err := q.GetFilecache(file)
	if err != nil {
		return "", nil, nil, err.
			Append(errors.LvlDebug, "Could not get file %v from the database", file.GetId()).
			AltStr(errors.LvlPlain, "Could not get file from the database")
	}
	return name, content, date, nil
}

func (file *DatabaseFile) GetContent(q types.DatabaseQueries) (io.Reader, *errors.ErrorTrace) {
	if file.content != nil {
		return bytes.NewReader(file.content), nil
	}

	name, reader, _, tr := file.fetchContentFromDatabase(q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not read contents of file %v", file.GetId()).
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not read file content from buffer").
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}

	file.name = name
	file.content = buf

	return bytes.NewReader(file.content), nil
}

func (file *DatabaseFile) GetBytes(q types.DatabaseQueries) ([]byte, *errors.ErrorTrace) {
	content, tr := file.GetContent(q)
	if tr != nil {
		return nil, tr
	}

	bytes, err := io.ReadAll(content)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not read file content from buffer").
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}

	return bytes, nil
}
