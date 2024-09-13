package config

import (
	"luna-backend/common"
	"luna-backend/db"

	"github.com/sirupsen/logrus"
)

type Api struct {
	Db           *db.Database
	CommonConfig *common.CommonConfig
	Logger       *logrus.Entry
	run          func(*Api)
}

func NewApi(db *db.Database, commonConfig *common.CommonConfig, logger *logrus.Entry, run func(*Api)) *Api {
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
