package queries

import (
	"bytes"
	"fmt"
	"io"
	"luna-backend/types"
	"time"
)

func (q *Queries) GetFilecache(file types.File) (io.Reader, *time.Time, error) {
	var content []byte
	var date time.Time

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT file, date
		FROM filecache
		WHERE id = $1;
		`,
		file.GetId().UUID(),
	).Scan(&content, &date)

	// TODO: read directly from the database instead of into an array first
	return bytes.NewReader(content), &date, err
}

func (q *Queries) SetFilecache(file types.File, content io.Reader) error {
	buf, err := io.ReadAll(content)
	if err != nil {
		return err
	}

	_, err = q.Tx.Exec(
		q.Context,
		`
		INSERT INTO filecache (id, file, date)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE
		SET file = $2;
		`,
		file.GetId().UUID(),
		buf,
	)

	return err
}

func (q *Queries) SetFilecacheWithoutId(file types.File, content io.Reader) (types.ID, error) {
	buf, err := io.ReadAll(content)
	if err != nil {
		return types.EmptyId(), err
	}

	query := `
		INSERT INTO filecache (file, date)
		VALUES ($1, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE
		SET file = $1
		RETURNING id;
	`

	var id types.ID
	err = q.Tx.QueryRow(q.Context, query, buf).Scan(&id)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not set filecache without id: %v", err)
	}

	file.SetId(id)

	return id, err
}

func (q *Queries) DeleteFilecache(file types.File) error {
	_, err := q.Tx.Exec(
		q.Context,
		`
		DELETE FROM filecache
		WHERE id = $1;
		`,
		file.GetId().UUID(),
	)

	return err
}
