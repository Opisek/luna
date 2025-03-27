package config

import (
	"luna-backend/errors"
	"luna-backend/types"
)

type CommonConfig struct {
	Version  types.Version
	Env      *Environmental
	Settings *GlobalSettings
}

func (c *CommonConfig) LoggingVerbosity() int {
	if c.Settings == nil {
		return errors.LvlPlain
	} else {
		return c.Settings.LoggingVerbosity.Verbosity
	}
}
