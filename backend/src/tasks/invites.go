package tasks

import (
	"luna-backend/db"
	"luna-backend/errors"

	"github.com/sirupsen/logrus"
)

func DeleteExpiredRegistrationInvites(tx *db.Transaction, logger *logrus.Entry) *errors.ErrorTrace {
	return tx.Queries().DeleteExpiredInvites()
}
