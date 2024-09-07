package caldav

import (
	"luna-backend/auth"
	"luna-backend/sources"
	"net/url"

	"github.com/emersion/go-webdav/caldav"
)

func NewCaldavSource(name string, url *url.URL, auth auth.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id:   sources.NewRandomSourceId(),
		name: name,
		auth: auth,
		settings: &CaldavSettings{
			Url: url,
		},
	}
}

func PackCaldavSource(id sources.SourceId, name string, settings *CaldavSettings, auth auth.AuthMethod) *CaldavSource {
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
			source.settings.Url.String(),
		)

		if err != nil {
			return nil, err
		}
	}
	return source.client, nil
}
