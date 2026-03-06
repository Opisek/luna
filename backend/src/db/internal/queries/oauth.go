package queries

import (
	"fmt"
	"luna-backend/db/internal/util"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) InsertOauthClient(client *types.OauthClient) *errors.ErrorTrace {
	encryptionKey, tr := util.GetGlobalEncryptionKey(q.CommonConfig)
	if tr != nil {
		return tr.
			Append(errors.LvlWordy, "Could not insert oauth client")
	}

	query := `
		INSERT INTO oauth_clients (name, client_id, client_secret, base_url, scope)
		VALUES ($1, $2, PGP_SYM_ENCRYPT($3, $6), $4, $5)
		RETURNING id;
	`

	params := make([]any, 6)
	params[0] = client.Name
	params[1] = client.ClientId
	params[2] = client.ClientSecret
	params[3] = client.BaseUrl
	params[4] = client.Scope
	params[5] = encryptionKey

	err := q.Tx.
		QueryRow(
			q.Context,
			query,
			params...,
		).Scan(&client.Id)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not insert oauth client").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) GetOauthClientById(id types.ID) (*types.OauthClient, *errors.ErrorTrace) {
	decryptionKey, tr := util.GetGlobalDecryptionKey(q.CommonConfig)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not get oauth client %v", id).
			AltStr(errors.LvlWordy, "Could not get oauth client")
	}

	query := `
		SELECT id, name, client_id, PGP_SYM_DECRYPT(client_secret, $2), base_url, scope
		FROM oauth_clients
		WHERE id = $1;
	`

	fmt.Println("decryption key:", decryptionKey)

	client := &types.OauthClient{}
	var rawBaseUrl string

	err := q.Tx.QueryRow(
		q.Context,
		query,
		id.UUID(),
		decryptionKey,
	).Scan(&client.Id, &client.Name, &client.ClientId, &client.ClientSecret, &rawBaseUrl, &client.Scope)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Oauth client %v does not exist", id).
			Append(errors.LvlDebug, "Could not get oauth client %v", id).
			AltStr(errors.LvlPlain, "Oauth client not found")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Database encountered an error").
			Append(errors.LvlDebug, "Could not get oauth client %v", id).
			AltStr(errors.LvlWordy, "Could not get oauth client").
			Append(errors.LvlPlain, "Database error")
	}

	if rawBaseUrl == "" {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "The authorization URL is empty").
			Append(errors.LvlDebug, "Could not get oauth client %v", id).
			AltStr(errors.LvlWordy, "Could not get oauth client").
			Append(errors.LvlPlain, "Database error")
	}
	client.BaseUrl, err = types.NewUrl(rawBaseUrl)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse base URL").
			Append(errors.LvlDebug, "Could not get oauth client %v", id).
			AltStr(errors.LvlWordy, "Could not get oauth client").
			Append(errors.LvlPlain, "Database error")
	}

	return client, nil
}

func (q *Queries) GetOauthClients() ([]*types.OauthClient, *errors.ErrorTrace) {
	query := `
		SELECT id, name, client_id, authorization_url
		FROM oauth_clients;
	`

	rows, err := q.Tx.Query(q.Context, query)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get oauth clients").
			Append(errors.LvlPlain, "Database error")
	}
	defer rows.Close()

	clients := make([]*types.OauthClient, 0)
	for rows.Next() {
		client := &types.OauthClient{}
		var rawBaseUrl string
		err = rows.Scan(&client.Id, &client.Name, &client.ClientId, &rawBaseUrl)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan oauth client").
				Append(errors.LvlPlain, "Database error")
		}
		if rawBaseUrl == "" {
			return nil, errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "The authorization URL is empty").
				Append(errors.LvlDebug, "Could not get oauth client %v", client.Id).
				AltStr(errors.LvlWordy, "Could not get oauth clients").
				Append(errors.LvlPlain, "Database error")
		}
		client.BaseUrl, err = types.NewUrl(rawBaseUrl)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse authorization URL").
				Append(errors.LvlDebug, "Could not get oauth client %v", client.Id).
				AltStr(errors.LvlWordy, "Could not get oauth clients").
				Append(errors.LvlPlain, "Database error")
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func (q *Queries) UpdateOauthClient(client *types.OauthClient) *errors.ErrorTrace {
	var query string
	var params []any

	if client.ClientSecret == "" {
		query = `
			UPDATE oauth_clients
			SET name = $1, client_id = $2, base_url = $3, scope = $4
			WHERE id = $5;
		`
		params = make([]any, 5)
	} else {
		query = `
			UPDATE oauth_clients
			SET name = $1, client_id = $2, client_secret = PGP_SYM_ENCRYPT($6, $7), base_url = $3, scope = $4
			WHERE id = $5;
		`
		params = make([]any, 7)
	}

	params[0] = client.Name
	params[1] = client.ClientId
	params[2] = client.BaseUrl
	params[3] = client.Scope
	params[4] = client.Id.UUID()
	if client.ClientSecret != "" {
		encryptionKey, tr := util.GetGlobalEncryptionKey(q.CommonConfig)
		if tr != nil {
			return tr.
				Append(errors.LvlDebug, "Could not update oauth client %v", client.Id).
				AltStr(errors.LvlWordy, "Could not update oauth client")
		}

		params[5] = client.ClientSecret
		params[6] = encryptionKey
	}

	_, err := q.Tx.Exec(q.Context, query, params...)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not update oauth client").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) DeleteOauthClient(id types.ID) *errors.ErrorTrace {
	query := `
		DELETE FROM oauth_clients
		WHERE id = $1;
	`

	_, err := q.Tx.Exec(q.Context, query, id.UUID())
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not delete oauth client").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}
