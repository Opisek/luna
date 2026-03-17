package tasks

import (
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/errors"

	"github.com/sirupsen/logrus"
)

func DeleteExpiredOauthAuthorizationRequests(tx *db.Transaction, logger *logrus.Entry, config *config.CommonConfig) *errors.ErrorTrace {
	return tx.Queries().DeleteExpiredOauthAuthorizationRequests()
}
