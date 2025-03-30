package queries

import (
	"bytes"
	"io"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) GetFilecache(file types.File) (string, io.Reader, *time.Time, *errors.ErrorTrace) {
	var name string
	var content []byte
	var date time.Time

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT name, file, date
		FROM filecache
		WHERE id = $1;
		`,
		file.GetId().UUID(),
	).Scan(&name, &content, &date)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return "", nil, nil, errors.New().Status(http.StatusNotFound).
				Append(errors.LvlPlain, "File not found")
		default:
			return "", nil, nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlPlain, "Database error")
		}

	}

	// TODO: read directly from the database instead of into an array first
	return name, bytes.NewReader(content), &date, nil
}

func (q *Queries) SetFilecache(file types.File, content io.Reader) *errors.ErrorTrace {
	buf, err := io.ReadAll(content)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Could not read from buffer").
			Append(errors.LvlWordy, "Could not save file cache").
			Append(errors.LvlPlain, "Database error")
	}

	_, err = q.Tx.Exec(
		q.Context,
		`
		INSERT INTO filecache (id, file, name, date)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE
		SET file = $2, name = $3;
		`,
		file.GetId().UUID(),
		buf,
		file.GetName(q),
	)

	return errors.New().Status(http.StatusInternalServerError).
		AddErr(errors.LvlDebug, err).
		Append(errors.LvlWordy, "Could not save file cache").
		Append(errors.LvlPlain, "Database error")
}

func (q *Queries) SetFilecacheWithoutId(file types.File, content io.Reader) (types.ID, *errors.ErrorTrace) {
	buf, err := io.ReadAll(content)
	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Could not read from buffer").
			Append(errors.LvlWordy, "Could not save new file cache").
			Append(errors.LvlPlain, "Database error")
	}

	query := `
		INSERT INTO filecache (file, name, date)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE
		SET file = $1, name = $2
		RETURNING id;
	`

	var id types.ID
	err = q.Tx.QueryRow(q.Context, query, buf, file.GetName(q)).Scan(&id)
	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not save new file cache").
			Append(errors.LvlPlain, "Database error")
	}

	file.SetId(id)

	return id, nil
}

func (q *Queries) DeleteFilecache(file types.File) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		DELETE FROM filecache
		WHERE id = $1;
		`,
		file.GetId().UUID(),
	)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlPlain, "File not found")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Database error")
	}
}
