package queries

import (
	"fmt"
	"luna-backend/auth"
	"luna-backend/db/internal/parsing"
	"luna-backend/db/internal/util"
	"luna-backend/errors"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) GetSource(userId types.ID, sourceId types.ID) (primitives.Source, *errors.ErrorTrace) {
	decryptionKey, tr := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not get source %v", sourceId).
			AltStr(errors.LvlBroad, "Could not get source")
	}

	scanner := parsing.NewPgxScanner(q.PrimitivesParser, q)
	scanner.ScheduleSource()
	cols, params := scanner.Variables(3)

	query := fmt.Sprintf(
		`
		SELECT %s
		FROM sources
		WHERE id = $1 AND userid = $2;
		`,
		cols,
	)

	err := q.Tx.QueryRow(
		q.Context,
		query,
		sourceId.UUID(),
		userId.UUID(),
		decryptionKey,
	).Scan(params...)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Source %v for user %v not found", sourceId, userId).
			AltStr(errors.LvlPlain, "Source not found").
			AltStr(errors.LvlBroad, "Could not get source")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get source %v for user %v", sourceId, userId).
			AltStr(errors.LvlBroad, "Could not get source")
	}

	source, tr := scanner.GetSource()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not parse aource %v for user %v", sourceId, userId).
			AltStr(errors.LvlWordy, "Could not parse source").
			AltStr(errors.LvlBroad, "Could not get source")
	}

	return source, nil

}

func (q *Queries) GetSourcesByUser(userId types.ID) ([]primitives.Source, *errors.ErrorTrace) {
	decryptionKey, tr := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get sources")
	}

	scanner := parsing.NewPgxScanner(q.PrimitivesParser, q)
	scanner.ScheduleSource()
	cols, params := scanner.Variables(2)

	query := fmt.Sprintf(
		`
		SELECT %s
		FROM sources
		WHERE userid = $1;
		`,
		cols,
	)

	rows, err := q.Tx.Query(
		q.Context,
		query,
		userId.UUID(),
		decryptionKey,
	)
	switch err {
	case nil:
		break
	case pgx.ErrNoRows: // I don't think this is actually possible
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Sources for user %v not found", userId).
			AltStr(errors.LvlPlain, "Sources not found").
			AltStr(errors.LvlBroad, "Could not get sources")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get sources for user %v", userId).
			AltStr(errors.LvlBroad, "Could not get sources")
	}
	defer rows.Close()

	sources := []primitives.Source{}
	for rows.Next() {
		rows.Scan(params...) // TODO: we might have to "reset" the scanner each time due to pass by reference
		source, err := scanner.GetSource()
		if err != nil {
			return nil, err.
				Append(errors.LvlDebug, "Could not parse sources for user %v", userId).
				AltStr(errors.LvlWordy, "Could not parse sources").
				AltStr(errors.LvlBroad, "Could not get sources")
		}
		sources = append(sources, source)
	}

	return sources, nil
}

// This is only used to refetch iCal files cache periodically.
// The information about file URL could be stored in another table instead,
// so we don't have to query the more sensitive sources table.
func (q *Queries) GetSourceSettingsByType(sourceType string) ([][]byte, *errors.ErrorTrace) {
	var err error

	rows, err := q.Tx.Query(
		q.Context,
		`
		SELECT settings
		FROM sources
		WHERE type = $1;
		`,
		sourceType,
	)
	switch err {
	case nil:
		break
	case pgx.ErrNoRows: // I don't think this is actually possible
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Sources of type %v not found", sourceType).
			AltStr(errors.LvlPlain, "Sources not found").
			AltStr(errors.LvlBroad, "Could not get sources")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get sources of type %v", sourceType).
			AltStr(errors.LvlBroad, "Could not get sources")
	}
	defer rows.Close()

	settings := [][]byte{}
	for rows.Next() {
		var setting []byte
		err = rows.Scan(&setting)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				AltStr(errors.LvlBroad, "Could not get sources")
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

func (q *Queries) InsertSource(userId types.ID, source primitives.Source) (types.ID, *errors.ErrorTrace) {
	encryptionKey, tr := util.GetUserEncryptionKey(q.CommonConfig, userId)
	if tr != nil {
		return types.EmptyId(), tr.
			Append(errors.LvlDebug, "Could not insert source %v for user %v", source.GetName(), userId).
			AltStr(errors.LvlWordy, "Could not insert source %v", source.GetName()).
			AltStr(errors.LvlPlain, "Could not add source %v", source.GetName()).
			AltStr(errors.LvlBroad, "Could not add source")
	}

	query := `
		INSERT INTO sources (userid, name, type, settings, auth_type, auth)
		VALUES ($1, $2, $3, $4, PGP_SYM_ENCRYPT($5, $7), PGP_SYM_ENCRYPT($6, $7))
		RETURNING id;
	`
	marshalledAuth, err := source.GetAuth().String()
	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not marshal authentication").
			Append(errors.LvlDebug, "Could not insert source %v for user %v", source.GetName(), userId).
			AltStr(errors.LvlWordy, "Could not insert source %v", source.GetName()).
			AltStr(errors.LvlPlain, "Could not add source %v", source.GetName()).
			AltStr(errors.LvlBroad, "Could not add source")
	}
	args := []any{userId.UUID(), source.GetName(), source.GetType(), source.GetSettings(), source.GetAuth().GetType(), marshalledAuth, encryptionKey}

	var id uuid.UUID
	err = q.Tx.QueryRow(q.Context, query, args...).Scan(&id)

	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not insert source %v for user %v", source.GetName(), userId).
			AltStr(errors.LvlWordy, "Could not insert source %v", source.GetName()).
			AltStr(errors.LvlPlain, "Could not add source %v", source.GetName()).
			AltStr(errors.LvlBroad, "Could not add source")
	}

	return types.IdFromUuid(id), nil
}

func (q *Queries) UpdateSource(userId types.ID, sourceId types.ID, newName string, newAuth auth.AuthMethod, newSourceType string, newSourceSettings primitives.SourceSettings) *errors.ErrorTrace {
	encryptionKey, tr := util.GetUserEncryptionKey(q.CommonConfig, userId)
	if tr != nil {
		return tr.
			Append(errors.LvlDebug, "Could not update source %v", sourceId).
			AltStr(errors.LvlWordy, "Could not update source").
			AltStr(errors.LvlBroad, "Could not edit source")
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
			return errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not marshal authentication").
				Append(errors.LvlDebug, "Could not update source %v", sourceId).
				AltStr(errors.LvlWordy, "Could not update source").
				AltStr(errors.LvlBroad, "Could not edit source")
		}
		args = append(args, newAuth.GetType(), marshalledAuth, encryptionKey)
	}

	if len(changes) == 0 {
		return errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Nothing to update").
			AltStr(errors.LvlPlain, "Nothing to change")
	}

	query := fmt.Sprintf(`
		UPDATE sources
		SET %s
		WHERE userid = $%d AND id = $%d;
	`, strings.Join(changes, ", "), len(args)+1, len(args)+2)
	args = append(args, userId.UUID(), sourceId.UUID())

	_, err := q.Tx.Exec(q.Context, query, args...)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Source %v not found", sourceId).
			AltStr(errors.LvlPlain, "Source not found").
			Append(errors.LvlDebug, "Could not update source %v", sourceId).
			Append(errors.LvlWordy, "Could not update source").
			Append(errors.LvlPlain, "Could not edit source")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update source %v", sourceId).
			Append(errors.LvlWordy, "Could not update source").
			Append(errors.LvlPlain, "Could not edit source")
	}
}

func (q *Queries) DeleteSource(userId types.ID, sourceId types.ID) (bool, *errors.ErrorTrace) {
	tag, err := q.Tx.Exec(
		q.Context,
		`
		DELETE FROM sources
		WHERE userid = $1 AND id = $2;
		`,
		userId.UUID(),
		sourceId,
	)

	switch err {
	case nil:
		return tag.RowsAffected() != 0, nil
	case pgx.ErrNoRows:
		return false, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Source %v not found", sourceId).
			AltStr(errors.LvlPlain, "Source not found").
			Append(errors.LvlDebug, "Could not delete source %v", sourceId).
			AltStr(errors.LvlBroad, "Could not delete source")
	default:
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not delete source %v", sourceId).
			AltStr(errors.LvlBroad, "Could not delete source")
	}
}
