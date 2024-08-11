package ical

import (
	"luna-backend/sources"
	"net/url"
)

func NewIcalSource(url *url.URL, auth sources.SourceAuth) *IcalSource {
	return &IcalSource{
		settings: &IcalSettings{
			Url:  url,
			Auth: auth,
		},
	}
}
