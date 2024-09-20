package parsing

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/types"
)

func ParseSource(id types.ID, name string, sourceType string, sourceSettings []byte, authType string, authBytes []byte) (primitives.Source, error) {
	var err error

	var authMethod auth.AuthMethod
	switch authType {
	case types.AuthNone:
		authMethod = auth.NewNoAuth()
	case types.AuthBasic:
		basicAuth := &auth.BasicAuth{}
		err = json.Unmarshal([]byte(authBytes), basicAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal basic auth: %v", err)
		}
		authMethod = basicAuth
	case types.AuthBearer:
		bearerAuth := &auth.BearerAuth{}
		err = json.Unmarshal([]byte(authBytes), bearerAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal bearer auth: %v", err)
		}
		authMethod = bearerAuth
	default:
		return nil, fmt.Errorf("unknown auth type: %v", authType)
	}

	switch sourceType {
	case types.SourceCaldav:
		settings := &caldav.CaldavSourceSettings{}
		err = json.Unmarshal(sourceSettings, settings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		caldavSource := caldav.PackCaldavSource(
			id,
			name,
			settings,
			authMethod,
		)
		return caldavSource, nil
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type: %v", sourceType)
	}
}
