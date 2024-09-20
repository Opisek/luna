package tables

import "context"

func (q *Tables) InitalizeVersionTable() error {
	// Keeps track of the current backend version as well as stores past
	// versions in case some specific migration rules need to be followed

	// Version table:
	// id major minor patch extension installed

	_, err := q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS version (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			major INT NOT NULL,
			minor INT NOT NULL,
			patch INT NOT NULL,
			extension VARCHAR(255),
			installed TIMESTAMP NOT NULL
		);
		`,
	)

	return err
}
