package files

import (
	"bytes"
	"fmt"
	"io"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/net"
	"luna-backend/types"
	"time"
)

// Implements types.File
type RemoteFile struct {
	url     *types.Url
	date    *time.Time
	content []byte
}

func NewRemoteFile(url *types.Url) *RemoteFile {
	return &RemoteFile{url: url}
}

func (file *RemoteFile) GetId() types.ID {
	return crypto.DeriveID(types.UrlNamespace(), file.url.URL().String())
}

func (file *RemoteFile) fetchContentFromRemote(q types.FileQueries) (io.Reader, error) {
	content, err := net.FetchFile(file.url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch remote file content: %w", err)
	}

	// TODO: don't use local buffer, instead run the database query in a
	// separate goroutine and return a reader directly without buffering the
	// whole file first
	var buf bytes.Buffer
	cache := io.TeeReader(content, &buf)
	err = q.SetFilecache(file, cache)
	if err != nil {
		// TODO: Logger.Warnf("could not set remote file cache in database: %v", err)
	}

	return &buf, nil
}

func (file *RemoteFile) fetchContentFromDatabase(q types.FileQueries) (io.Reader, *time.Time, error) {
	content, date, err := q.GetFilecache(file)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get file content from database: %w", err)
	}
	return content, date, nil
}

func (file *RemoteFile) GetContent(q types.FileQueries) (io.Reader, error) {
	curTime := time.Now()

	var err error
	var reader io.Reader

	// Try to get from the database first
	if file.content == nil {
		reader, file.date, _ = file.fetchContentFromDatabase(q)
	}

	// Refetch from the remote based on the cache lifetime
	if reader != nil {
		deltaTime := curTime.Sub(*file.date)

		if deltaTime >= common.LifetimeCacheSoft {
			reader, err = file.fetchContentFromRemote(q)

			if err == nil {
				file.date = &curTime
			} else if deltaTime >= common.LifetimeCacheHard {
				return nil, fmt.Errorf("could not get file content: %w", err)
			}
		}
	} else {
		reader, err = file.fetchContentFromRemote(q)
		if err != nil {
			return nil, fmt.Errorf("could not get file content: %w", err)
		}
		file.date = &curTime
	}

	// TODO: figure out a proper way to use a reader without ending up saving the whole content to an array in the process
	file.content, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read file content: %w", err)
	}
	return bytes.NewReader(file.content), nil
}
