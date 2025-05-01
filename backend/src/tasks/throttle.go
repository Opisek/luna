package tasks

import (
	"luna-backend/db"
	"luna-backend/errors"

	"github.com/sirupsen/logrus"
)

type ThrottleInterface interface {
	CleanStaleEntries()
}

func DeleteStaleRequestThrottleEntries(throttle ThrottleInterface) func(tx *db.Transaction, logger *logrus.Entry) *errors.ErrorTrace {
	return func(tx *db.Transaction, logger *logrus.Entry) *errors.ErrorTrace {
		throttle.CleanStaleEntries()
		return nil
	}
}
