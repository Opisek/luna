package caldav

import (
	"luna-backend/auth"
	"luna-backend/types"

	"github.com/emersion/go-webdav/caldav"
)

func NewCaldavSource(name string, url *types.Url, auth auth.AuthMethod) *CaldavSource {
	return &CaldavSource{
		//id:   types.RandomId(),
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth,
		settings: &CaldavSettings{
			Url: url,
		},
	}
}

func PackCaldavSource(id types.ID, name string, settings *CaldavSettings, auth auth.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
	}
}

func (source *CaldavSource) getClient() (*caldav.Client, error) {
	if source.client == nil {
		var err error
		source.client, err = caldav.NewClient(
			source.auth,
			source.settings.Url.URL().String(),
		)

		if err != nil {
			return nil, err
		}
	}
	return source.client, nil
}
