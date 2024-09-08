package ical

import (
	"luna-backend/auth"
	"luna-backend/types"
)

func NewIcalSource(url *types.Url, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		//id: types.RandomId(),
		id: types.EmptyId(), // Placeholder until the database assigns an ID
		settings: &IcalSettings{
			Url:  url,
			Auth: auth,
		},
	}
}
