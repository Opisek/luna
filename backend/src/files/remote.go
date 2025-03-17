package files

import (
	"bytes"
	"io"
	"luna-backend/auth"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/net"
	"luna-backend/types"
	"net/http"
	"path"
	"time"
)

// Implements types.File
type RemoteFile struct {
	url     *types.Url
	date    *time.Time
	content []byte
	auth    auth.AuthMethod
}

func NewRemoteFile(url *types.Url, auth auth.AuthMethod) *RemoteFile {
	return &RemoteFile{url: url, auth: auth}
}

func (file *RemoteFile) GetId() types.ID {
	return crypto.DeriveID(types.UrlNamespace(), file.url.String())
}

func (file *RemoteFile) SetId(id types.ID) {
	panic("illegal operation")
}

func (file *RemoteFile) GetName(_ types.DatabaseQueries) string {
	return path.Base(file.url.URL().Path)
}

func (file *RemoteFile) fetchContentFromRemote(q types.DatabaseQueries) (io.Reader, *errors.ErrorTrace) {
	content, err := net.FetchFile(file.url, file.auth, q.GetContext())
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not read from remote").
			Append(errors.LvlPlain, "Could not read contents of file")
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

func (file *RemoteFile) fetchContentFromDatabase(q types.DatabaseQueries) (io.Reader, *time.Time, *errors.ErrorTrace) {
	_, content, date, err := q.GetFilecache(file)
	if err != nil {
		return nil, nil, err.
			Append(errors.LvlDebug, "Could not get file %v from the database", file.GetId()).
			AltStr(errors.LvlPlain, "Could not get file from the database")
	}
	return content, date, nil
}

func (file *RemoteFile) GetContent(q types.DatabaseQueries) (io.Reader, *errors.ErrorTrace) {
	curTime := time.Now()

	var tr *errors.ErrorTrace
	var reader io.Reader

	// Try to get from the database first
	if file.content == nil {
		reader, file.date, _ = file.fetchContentFromDatabase(q)
	}
	if file.content != nil {
		reader = bytes.NewReader(file.content)
	}

	// Refetch from the remote based on the cache lifetime
	if reader != nil {
		deltaTime := curTime.Sub(*file.date)

		if deltaTime >= common.LifetimeCacheSoft {
			reader, tr = file.fetchContentFromRemote(q)

			if tr == nil {
				file.date = &curTime
			} else if deltaTime >= common.LifetimeCacheHard {
				return nil, tr
			}
		}
	} else {
		reader, tr = file.fetchContentFromRemote(q)
		if tr != nil {
			return nil, tr
		}
		file.date = &curTime
	}

	// TODO: figure out a proper way to use a reader without ending up saving the whole content to an array in the process
	var err error
	file.content, err = io.ReadAll(reader)
	if tr != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not read from buffer").
			Append(errors.LvlDebug, "Could not read contents of file %v at %v", file.GetId(), file.url).
			AltStr(errors.LvlWordy, "Could not read contents of file at %v", file.url).
			AltStr(errors.LvlPlain, "Could not read contents of file")
	}
	return bytes.NewReader(file.content), nil
}

func (file *RemoteFile) ForceFetchFromRemote(q types.DatabaseQueries) *errors.ErrorTrace {
	_, err := file.fetchContentFromRemote(q)
	return err
}

func (file *RemoteFile) GetBytes(q types.DatabaseQueries) ([]byte, *errors.ErrorTrace) {
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
