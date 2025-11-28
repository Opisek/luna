package queries

import (
	"luna-backend/errors"
	"luna-backend/perms"
	"luna-backend/types"
)

// TODO: check if the session belongs to the user making the request?

func (q *Queries) GetTokenPermissions(sessionid types.ID) (*perms.TokenPermissions, *errors.ErrorTrace) {
	rows, err := q.Tx.Query(
		q.Context,
		`
		SELECT permission
		FROM token_permissions
		WHERE sessionid = $1;
		`,
		sessionid,
	)
	if err != nil {
		return nil, errors.New().Status(500).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get token permissions").
			AltStr(errors.LvlPlain, "Database error")
	}
	defer rows.Close()

	permissions := []perms.Permission{}
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, errors.New().Status(500).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not read token permission").
				AltStr(errors.LvlPlain, "Database error")
		}
		permissions = append(permissions, perms.Permission(permission))
	}

	return perms.FromList(permissions), nil
}

func (q *Queries) UpdateTokenPermissions(sessionid types.ID, permissions *perms.TokenPermissions) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		DELETE FROM token_permissions
		WHERE sessionid = $1;
		`,
		sessionid,
	)
	if err != nil {
		return errors.New().Status(500).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not clear existing token permissions").
			AltStr(errors.LvlPlain, "Database error")
	}

	// TODO: bulk insert?
	for _, permission := range permissions.ToList() {
		_, err := q.Tx.Exec(
			q.Context,
			`
			INSERT INTO token_permissions (sessionid, permission)
			VALUES ($1, $2);
			`,
			sessionid,
			string(permission),
		)
		if err != nil {
			return errors.New().Status(500).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not set token permission").
				AltStr(errors.LvlPlain, "Database error")
		}
	}

	return nil
}
