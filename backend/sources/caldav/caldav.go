package caldav

import (
	"net/http"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

func NewCaldavSource(settings *CaldavSettings) *CaldavSource {
	return &CaldavSource{
		settings: settings,
	}
}

func (source *CaldavSource) getClient() (*caldav.Client, error) {
	if source.client == nil {
		var err error
		source.client, err = caldav.NewClient(
			webdav.HTTPClientWithBasicAuth(
				http.DefaultClient,
				source.settings.Username,
				source.settings.Password,
			),
			source.settings.Url.String(),
		)
		if err != nil {
			return nil, err
		}
	}
	return source.client, nil
}
