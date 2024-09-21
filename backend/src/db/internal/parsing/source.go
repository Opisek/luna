package parsing

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/db/internal/tables"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/types"
)

func ParseSource(entry *tables.SourceEntry) (primitives.Source, error) {
	var err error

	var authMethod auth.AuthMethod
	switch entry.AuthType {
	case types.AuthNone:
		authMethod = auth.NewNoAuth()
	case types.AuthBasic:
		basicAuth := &auth.BasicAuth{}
		err = json.Unmarshal([]byte(entry.Auth), basicAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal basic auth: %v", err)
		}
		authMethod = basicAuth
	case types.AuthBearer:
		bearerAuth := &auth.BearerAuth{}
		err = json.Unmarshal([]byte(entry.Auth), bearerAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal bearer auth: %v", err)
		}
		authMethod = bearerAuth
	default:
		return nil, fmt.Errorf("unknown auth type: %v", entry.Auth)
	}

	switch entry.Type {
	case types.SourceCaldav:
		settings := &caldav.CaldavSourceSettings{}
		err = json.Unmarshal(entry.Settings, settings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		caldavSource := caldav.PackCaldavSource(
			entry.Id,
			entry.Name,
			settings,
			authMethod,
		)
		return caldavSource, nil
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type: %v", entry.Type)
	}
}
