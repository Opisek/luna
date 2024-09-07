package caldav

import (
	"luna-backend/auth"
	"luna-backend/sources"
	"net/url"

	"github.com/emersion/go-webdav/caldav"
)

func NewCaldavSource(url *url.URL, auth auth.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id: sources.NewRandomSourceId(),
		settings: &CaldavSettings{
			Url:  url,
			Auth: auth,
		},
	}
}

func (source *CaldavSource) getClient() (*caldav.Client, error) {
	if source.client == nil {
		var err error
		source.client, err = caldav.NewClient(
			source.settings.Auth,
			source.settings.Url.String(),
		)

		if err != nil {
			return nil, err
		}
	}
	return source.client, nil
}
