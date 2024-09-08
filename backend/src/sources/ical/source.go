package ical

import (
	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/types"
)

func NewIcalSource(url *types.Url, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		id: sources.NewRandomSourceId(),
		settings: &IcalSettings{
			Url:  url,
			Auth: auth,
		},
	}
}
