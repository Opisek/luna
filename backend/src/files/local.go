package files

import (
	"bytes"
	"fmt"
	"io"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"os"
	"time"
)

// Implements types.File
type LocalFile struct {
	path    *types.Path
	date    *time.Time
	content []byte
}

func NewLocalFile(path *types.Path) *LocalFile {
	return &LocalFile{path: path}
}

func (file *LocalFile) GetId() types.ID {
	return crypto.DeriveID(types.UrlNamespace(), file.path.String())
}

func (file *LocalFile) SetId(id types.ID) {
	panic("illegal operation")
}

func (file *LocalFile) fetchContentFromFilesystem() (io.Reader, *errors.ErrorTrace) {
	fd, err := os.Open(file.path.String())
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not open file %v at %v", file.GetId(), file.path).
			AltStr(errors.LvlWordy, "Could not open file at %v", file.path).
			AltStr(errors.LvlPlain, "Could not open file")
	}

	defer func() {
		if err := fd.Close(); err != nil {
			panic(fmt.Errorf("could not close file: %w", err))
		}
	}()

	buf, err := io.ReadAll(fd)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not read from filesystem").
			Append(errors.LvlDebug, "Could not read contents of file %v at %v", file.GetId(), file.path).
			AltStr(errors.LvlWordy, "Could not read contents of file at %v", file.path).
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}

	return bytes.NewReader(buf), nil
}

func (file *LocalFile) GetContent(q types.DatabaseQueries) (io.Reader, *errors.ErrorTrace) {
	curTime := time.Now()

	var tr *errors.ErrorTrace
	var reader io.Reader

	if file.content == nil {
		reader, tr = file.fetchContentFromFilesystem()
		if tr != nil {
			return nil, tr
		}
		file.date = &curTime
	}
	if file.content != nil {
		reader = bytes.NewReader(file.content)
	}

	deltaTime := curTime.Sub(*file.date)

	if deltaTime >= common.LifetimeCacheSoft {
		reader, tr = file.fetchContentFromFilesystem()

		if tr == nil {
			file.date = &curTime
		} else if deltaTime >= common.LifetimeCacheHard {
			return nil, tr
		}
	}

	// TODO: figure out a proper way to use a reader without ending up saving the whole content to an array in the process
	var err error
	file.content, err = io.ReadAll(reader)
	if tr != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not read from buffer").
			Append(errors.LvlDebug, "Could not read contents of file %v at %v", file.GetId(), file.path).
			AltStr(errors.LvlWordy, "Could not read contents of file at %v", file.path).
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}
	return bytes.NewReader(file.content), nil
}
