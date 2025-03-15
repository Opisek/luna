package files

import (
	"bytes"
	"fmt"
	"io"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/types"
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

func (file *LocalFile) fetchContentFromFilesystem() (io.Reader, error) {
	fd, err := os.Open(file.path.String())
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	defer func() {
		if err := fd.Close(); err != nil {
			panic(fmt.Errorf("could not close file: %w", err))
		}
	}()

	buf, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return bytes.NewReader(buf), nil
}

func (file *LocalFile) GetContent(q types.DatabaseQueries) (io.Reader, error) {
	curTime := time.Now()

	var err error
	var reader io.Reader

	if file.content == nil {
		reader, err = file.fetchContentFromFilesystem()
		if err != nil {
			return nil, fmt.Errorf("could not get file content: %w", err)
		}
		file.date = &curTime
	}
	if file.content != nil {
		reader = bytes.NewReader(file.content)
	}

	deltaTime := curTime.Sub(*file.date)

	if deltaTime >= common.LifetimeCacheSoft {
		reader, err = file.fetchContentFromFilesystem()

		if err == nil {
			file.date = &curTime
		} else if deltaTime >= common.LifetimeCacheHard {
			return nil, fmt.Errorf("could not get file content: %w", err)
		}
	}

	// TODO: figure out a proper way to use a reader without ending up saving the whole content to an array in the process
	file.content, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read file content: %w", err)
	}
	return bytes.NewReader(file.content), nil
}
