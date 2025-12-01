package tasks

import (
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/errors"

	"github.com/sirupsen/logrus"
)

type ThrottleInterface interface {
	CleanStaleEntries()
}

func DeleteStaleRequestThrottleEntries(throttle ThrottleInterface) func(tx *db.Transaction, logger *logrus.Entry, config *config.CommonConfig) *errors.ErrorTrace {
	return func(tx *db.Transaction, logger *logrus.Entry, config *config.CommonConfig) *errors.ErrorTrace {
		throttle.CleanStaleEntries()
		return nil
	}
}
