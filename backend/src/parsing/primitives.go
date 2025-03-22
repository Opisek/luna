package parsing

import (
	"encoding/json"
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/protocols/caldav"
	"luna-backend/protocols/ical"
	"luna-backend/types"
	"net/http"
)

type PrimitivesParser struct{}

func GetPrimitivesParser() PrimitivesParser {
	return PrimitivesParser{}
}

func (PrimitivesParser) ParseSource(entry *types.SourceDatabaseEntry) (types.Source, *errors.ErrorTrace) {
	var err error

	var authMethod types.AuthMethod
	switch entry.AuthType {
	case types.AuthNone:
		authMethod = auth.NewNoAuth()
	case types.AuthBasic:
		basicAuth := &auth.BasicAuth{}
		err = json.Unmarshal([]byte(entry.Auth), basicAuth)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal basic authentication").
				Append(errors.LvlWordy, "Could not unmarshal authentication")
		}
		authMethod = basicAuth
	case types.AuthBearer:
		bearerAuth := &auth.BearerAuth{}
		err = json.Unmarshal([]byte(entry.Auth), bearerAuth)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal bearer authentication").
				Append(errors.LvlWordy, "Could not unmarshal authentication")
		}
		authMethod = bearerAuth
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlPlain, "Unknown authentication type: %v", entry.Auth)
	}

	switch entry.Type {
	case types.SourceCaldav:
		settings := &caldav.CaldavSourceSettings{}
		err = json.Unmarshal(entry.Settings, settings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal CalDAV settings").
				Append(errors.LvlWordy, "Could not unmarshal settings")
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
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal iCal settings").
				Append(errors.LvlWordy, "Could not unmarshal settings")
		}
		icalSource, tr := ical.PackIcalSource(
			entry.Id,
			entry.Name,
			settings,
			authMethod,
		)
		if tr != nil {
			return nil, tr
		}
		return icalSource, nil
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Unknown source type: %v", entry.Type)
	}
}

func (PrimitivesParser) ParseCalendarSettings(sourceType string, settings []byte) (types.CalendarSettings, *errors.ErrorTrace) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavCalendarSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal CalDAV settings").
				Append(errors.LvlWordy, "Could not unmarshal settings")
		}
		return parsedSettings, nil
	case types.SourceIcal:
		parsedSettings := &ical.IcalCalendarSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal iCal settings").
				Append(errors.LvlWordy, "Could not unmarshal settings")
		}
		return parsedSettings, nil
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Unknown source type: %v", sourceType)
	}
}

func (PrimitivesParser) ParseEventSettings(sourceType string, settings []byte) (types.EventSettings, *errors.ErrorTrace) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavEventSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal CalDAV settings").
				Append(errors.LvlWordy, "Could not unmarshal settings")
		}
		return parsedSettings, nil
	case types.SourceIcal:
		parsedSettings := &ical.IcalEventSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not unmarshal iCal settings").
				Append(errors.LvlWordy, "Could not unmarshal settings")
		}
		return parsedSettings, nil
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Unknown source type: %v", sourceType)
	}
}
