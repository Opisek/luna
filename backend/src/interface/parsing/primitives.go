package parsing

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/interface/protocols/ical"
	"luna-backend/types"
)

type PrimitivesParser struct{}

func GetPrimitivesParser() PrimitivesParser {
	return PrimitivesParser{}
}

func (PrimitivesParser) ParseSource(entry *types.SourceDatabaseEntry) (primitives.Source, error) {
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
		settings := &ical.IcalSourceSettings{}
		err = json.Unmarshal(entry.Settings, settings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal ical settings: %v", err)
		}
		icalSource := ical.PackIcalSource(
			entry.Id,
			entry.Name,
			settings,
			authMethod,
		)
		return icalSource, nil
	default:
		return nil, fmt.Errorf("unknown source type: %v", entry.Type)
	}
}

func (PrimitivesParser) ParseCalendarSettings(sourceType string, settings []byte) (primitives.CalendarSettings, error) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavCalendarSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		return parsedSettings, nil
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type %v", sourceType)
	}
}

func (PrimitivesParser) ParseEventSettings(sourceType string, settings []byte) (primitives.EventSettings, error) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavEventSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		return parsedSettings, nil
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type %v", sourceType)
	}
}
