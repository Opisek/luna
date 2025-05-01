package queries

import (
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) InsertInvite(invite *types.RegistrationInvite) *errors.ErrorTrace {
	query := `
		INSERT INTO invites (author, email, expires_at, code)
		VALUES ($1, $2, $3, $4)
		RETURNING inviteid;
	`

	params := make([]any, 4)
	params[0] = invite.Author.UUID()
	if invite.Email == "" {
		params[1] = nil
	} else {
		params[1] = invite.Email
	}
	params[2] = invite.Expires
	params[3] = invite.Code

	err := q.Tx.
		QueryRow(
			q.Context,
			query,
			params...,
		).Scan(&invite.InviteId)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not insert invite").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) GetValidInvites() ([]*types.RegistrationInvite, *errors.ErrorTrace) {
	query := `
		SELECT inviteid, author, COALESCE(email, ''), created_at, expires_at, code
		FROM invites
		WHERE expires_at > NOW()
		ORDER BY created_at DESC;
	`

	rows, err := q.Tx.Query(q.Context, query)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get invites").
			Append(errors.LvlPlain, "Database error")
	}
	defer rows.Close()

	invites := make([]*types.RegistrationInvite, 0)
	for rows.Next() {
		invite := &types.RegistrationInvite{}
		err = rows.Scan(&invite.InviteId, &invite.Author, &invite.Email, &invite.CreatedAt, &invite.Expires, &invite.Code)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan invite").
				Append(errors.LvlPlain, "Database error")
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

func (q *Queries) GetValidInvite(email string, code string) (*types.RegistrationInvite, *errors.ErrorTrace) {
	query := `
		SELECT inviteid, author, COALESCE(email, ''), created_at, expires_at, code
		FROM invites
		WHERE code = $1
		AND expires_at > NOW()
		AND (email IS NULL OR email = $2);
	`

	invite := &types.RegistrationInvite{}
	err := q.Tx.
		QueryRow(
			q.Context,
			query,
			code,
			email,
		).Scan(&invite.InviteId, &invite.Author, &invite.CreatedAt, &invite.Expires, &invite.Code)

	switch err {
	case nil:
		return invite, nil
	case pgx.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get invite").
			Append(errors.LvlPlain, "Database error")
	}
}

func (q *Queries) DeleteExpiredInvites() *errors.ErrorTrace {
	query := `
		DELETE FROM invites
		WHERE expires_at <= NOW();
	`

	_, err := q.Tx.Exec(q.Context, query)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not delete expired invites").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) DeleteInvite(inviteId types.ID) *errors.ErrorTrace {
	query := `
		DELETE FROM invites
		WHERE inviteid = $1;
	`

	_, err := q.Tx.Exec(q.Context, query, inviteId.UUID())
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not delete invite").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) DeleteInvites() *errors.ErrorTrace {
	query := `
		DELETE FROM invites;
	`

	_, err := q.Tx.Exec(q.Context, query)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not delete invites").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}
