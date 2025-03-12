package queries

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"luna-backend/types"
	"time"
)

func (q *Queries) GetFilecache(file types.File) (io.Reader, *time.Time, error) {
	var content []byte
	var date time.Time

	err := q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT content, date
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
		context.TODO(),
		`
		INSERT INTO filecache (id, content, date)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE
		SET content = $2;
		`,
		file.GetId().UUID(),
		buf,
	)

	fmt.Println(err)

	return err
}

func (q *Queries) DeleteFilecache(file types.File) error {
	_, err := q.Tx.Exec(
		context.TODO(),
		`
		DELETE FROM filecache
		WHERE id = $1;
		`,
		file.GetId().UUID(),
	)

	return err
}
