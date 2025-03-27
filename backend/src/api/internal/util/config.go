package util

import (
	"luna-backend/config"
	"luna-backend/db"

	"github.com/sirupsen/logrus"
)

type Api struct {
	Db           *db.Database
	CommonConfig *config.CommonConfig
	Logger       *logrus.Entry
	run          func(*Api)
}

func NewApi(db *db.Database, commonConfig *config.CommonConfig, logger *logrus.Entry, run func(*Api)) *Api {
	return &Api{
		Db:           db,
		CommonConfig: commonConfig,
		Logger:       logger,
		run:          run,
	}
}

func (api *Api) Run() {
	api.run(api)
}
