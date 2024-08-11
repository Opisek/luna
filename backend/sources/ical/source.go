package ical

import (
	"luna-backend/auth"
	"luna-backend/sources"
	"net/url"
)

func NewIcalSource(url *url.URL, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		id: sources.NewRandomSourceId(),
		settings: &IcalSettings{
			Url:  url,
			Auth: auth,
		},
	}
}
