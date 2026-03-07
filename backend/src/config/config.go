package config

import (
	"luna-backend/errors"
	"luna-backend/types"
)

type CommonConfig struct {
	Version                  types.Version
	Env                      *Environmental
	PublicUrl                *types.Url
	Settings                 *GlobalSettings
	TokenInvalidationChannel chan *types.Session
}

func (c *CommonConfig) LoggingVerbosity() int {
	if c.Settings == nil {
		return errors.LvlPlain
	} else {
		return c.Settings.LoggingVerbosity.Verbosity
	}
}
