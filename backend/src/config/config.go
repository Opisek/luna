package config

import (
	"luna-backend/cache"
	"luna-backend/errors"
	"luna-backend/types"
)

type CommonConfig struct {
	Version                  types.Version
	Env                      *Environmental
	Cache                    *cache.Cache
	PublicUrl                *types.Url
	Settings                 *GlobalSettings
	TokenInvalidationChannel chan *types.Session
	OauthInvalidationChannel chan types.ID
}

func (c *CommonConfig) LoggingVerbosity() int {
	if c.Settings == nil {
		return errors.LvlPlain
	} else {
		return c.Settings.LoggingVerbosity.Verbosity
	}
}
