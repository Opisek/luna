package queries

import (
	"context"
	"fmt"
	"luna-backend/auth"
	"luna-backend/db/internal/parsing"
	"luna-backend/db/internal/util"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"strings"

	"github.com/google/uuid"
)

func scanSource(rows types.PgxScanner) (primitives.Source, error) {
	var err error

	var id types.ID
	var name string
	var sourceType string
	var settings []byte
	var authType string
	var authBytes []byte
	err = rows.Scan(&id, &name, &sourceType, &settings, &authType, &authBytes)
	if err != nil {
		return nil, fmt.Errorf("could not scan source row: %v", err)
	}

	return parsing.ParseSource(id, name, sourceType, settings, authType, authBytes)
}

func (q *Queries) GetSource(userId types.ID, sourceId types.ID) (primitives.Source, error) {
	decryptionKey, err := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	row := q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT id, name, type, settings, PGP_SYM_DECRYPT(auth_type, $3), PGP_SYM_DECRYPT(auth, $3)
		FROM sources
		WHERE userid = $1 AND id = $2;
		`,
		userId.UUID(),
		sourceId.UUID(),
		decryptionKey,
	)

	source, err := scanSource(row)
	if err != nil {
		return nil, fmt.Errorf("could not parse source: %v", err)
	}

	return source, nil
}

func (q *Queries) GetSources(userId types.ID) ([]primitives.Source, error) {
	var err error

	decryptionKey, err := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	rows, err := q.Tx.Query(
		context.TODO(),
		`
		SELECT id, name, type, settings, PGP_SYM_DECRYPT(auth_type, $2), PGP_SYM_DECRYPT(auth, $2)
		FROM sources
		WHERE userid = $1;
		`,
		userId.UUID(),
		decryptionKey,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get sources: %v", err)
	}
	defer rows.Close()

	sources := []primitives.Source{}
	for rows.Next() {
		source, err := scanSource(rows)
		if err != nil {
			return nil, fmt.Errorf("could not parse source: %v", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (q *Queries) InsertSource(userId types.ID, source primitives.Source) (types.ID, error) {
	encryptionKey, err := util.GetUserEncryptionKey(q.CommonConfig, userId)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not get user encryption key: %v", err)
	}

	query := `
		INSERT INTO sources (userid, name, type, settings, auth_type, auth)
		VALUES ($1, $2, $3, $4, PGP_SYM_ENCRYPT($5, $7), PGP_SYM_ENCRYPT($6, $7))
		RETURNING id;
	`
	marshalledAuth, err := source.GetAuth().String()
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not marshal auth: %v", err)
	}
	args := []any{userId.UUID(), source.GetName(), source.GetType(), source.GetSettings(), source.GetAuth().GetType(), marshalledAuth, encryptionKey}

	var id uuid.UUID
	err = q.Tx.QueryRow(context.TODO(), query, args...).Scan(&id)

	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not insert source: %v", err)
	}

	return types.IdFromUuid(id), nil
}

func (q *Queries) UpdateSource(userId types.ID, sourceId types.ID, newName string, newAuth auth.AuthMethod, newSourceType string, newSourceSettings primitives.SourceSettings) error {
	encryptionKey, err := util.GetUserEncryptionKey(q.CommonConfig, userId)
	if err != nil {
		return fmt.Errorf("could not get user encryption key: %v", err)
	}

	changes := []string{}
	args := []any{}

	if newName != "" {
		changes = append(changes, fmt.Sprintf("name = $%d", len(changes)+1))
		args = append(args, newName)
	}
	if newSourceType != "" && newSourceSettings != nil {
		changes = append(changes, fmt.Sprintf("type = $%d", len(changes)+1), fmt.Sprintf("settings = $%d", len(changes)+2))
		args = append(args, newSourceType, newSourceSettings)
	}
	if newAuth != nil {
		changes = append(
			changes,
			fmt.Sprintf("auth_type = PGP_SYM_ENCRYPT($%d, $%d)", len(changes)+1, len(changes)+3),
			fmt.Sprintf("auth = PGP_SYM_ENCRYPT($%d, $%d)", len(changes)+2, len(changes)+3),
		)
		marshalledAuth, err := newAuth.String()
		if err != nil {
			return fmt.Errorf("could not marshal auth: %v", err)
		}
		args = append(args, newAuth.GetType(), marshalledAuth, encryptionKey)
	}

	if len(changes) == 0 {
		return fmt.Errorf("no changes to update")
	}

	query := fmt.Sprintf(`
		UPDATE sources
		SET %s
		WHERE userid = $%d AND id = $%d;
	`, strings.Join(changes, ", "), len(args)+1, len(args)+2)
	args = append(args, userId.UUID(), sourceId.UUID())

	_, err = q.Tx.Exec(context.TODO(), query, args...)

	if err != nil {
		return fmt.Errorf("could not update source: %v", err)
	}

	return nil
}

func (q *Queries) DeleteSource(userId types.ID, sourceId types.ID) (bool, error) {
	tag, err := q.Tx.Exec(
		context.TODO(),
		`
		DELETE FROM sources
		WHERE userid = $1 AND id = $2;
		`,
		userId.UUID(),
		sourceId,
	)
	if err != nil {
		return false, fmt.Errorf("could not delete source: %v", err)
	}
	return tag.RowsAffected() != 0, nil
}
