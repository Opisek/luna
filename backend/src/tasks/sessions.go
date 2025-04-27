package tasks

import (
	"luna-backend/db"
	"luna-backend/errors"
	"time"

	"github.com/sirupsen/logrus"
)

func DeleteStaleShortLivedSessions(tx *db.Transaction, logger *logrus.Entry) *errors.ErrorTrace {
	currentTime := time.Now()

	return tx.Queries().DeleteExpiredSessions(currentTime.Add(-time.Hour), true)
}

func DeleteStaleLongLivedSessions(tx *db.Transaction, logger *logrus.Entry) *errors.ErrorTrace {
	currentTime := time.Now()

	return tx.Queries().DeleteExpiredSessions(currentTime.AddDate(0, -1, 0), true)
}
