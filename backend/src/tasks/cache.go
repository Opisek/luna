package tasks

import (
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/errors"

	"github.com/sirupsen/logrus"
)

func ClearStaleCache(_ *db.Transaction, _ *logrus.Entry, config *config.CommonConfig) *errors.ErrorTrace {
	config.Cache.DeleteStaleEntries()
	return nil
}
