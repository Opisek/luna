package queries

import (
	"context"
	"fmt"
	"luna-backend/types"
)

func (q *Queries) GetPassword(id types.ID) (*types.PasswordEntry, error) {
	var err error

	entry := &types.PasswordEntry{}

	err = q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT hash, salt, algorithm, parameters
		FROM passwords
		WHERE userid = $1;
		`, id.UUID(),
	).Scan(&entry.Hash, &entry.Salt, &entry.Algorithm, &entry.Parameters)

	if err != nil {
		return nil, fmt.Errorf("could not get password information of user %v: %v", id, err)
	}
	return entry, nil
}

func (q *Queries) InsertPassword(userId types.ID, entry *types.PasswordEntry) error {
	var err error

	_, err = q.Tx.Exec(
		context.TODO(),
		`
		INSERT INTO passwords (userid, hash, salt, algorithm, parameters)
		VALUES ($1, $2, $3, $4, $5);
		`,
		userId.UUID(), entry.Hash, entry.Salt, entry.Algorithm, entry.Parameters,
	)

	if err != nil {
		return fmt.Errorf("could not insert password information of user %v: %v", userId, err)
	}
	return err
}

func (q *Queries) UpdatePassword(userId types.ID, entry *types.PasswordEntry) error {
	var err error

	_, err = q.Tx.Exec(
		context.TODO(),
		`
		UPDATE passwords
		SET hash = $1, salt = $2, algorithm = $3, parameters = $4
		WHERE userid = $5;
		`,
		entry.Hash, entry.Salt, entry.Algorithm, entry.Parameters,
		userId.UUID(),
	)

	if err != nil {
		return fmt.Errorf("could not update password information of user %v: %v", userId, err)
	}
	return err
}
