package services

import (
	"context"
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/errors"
	"luna-backend/types"
	"time"

	"github.com/sirupsen/logrus"
)

type TokenInvalidationService struct {
	receiveChannel chan *types.Session
	db             *db.Database
	commonConfig   *config.CommonConfig
	logger         *logrus.Entry
}

func NewTokenInvalidationService(db *db.Database, commonConfig *config.CommonConfig, logger *logrus.Entry) *TokenInvalidationService {
	service := TokenInvalidationService{
		receiveChannel: make(chan *types.Session),
		db:             db,
		commonConfig:   commonConfig,
		logger:         logger,
	}

	commonConfig.TokenInvalidationChannel = service.Channel()

	return &service
}

func (t *TokenInvalidationService) Start() {
	t.logger.Info("starting token invalidation service")
	go func() {
		for s := range t.receiveChannel {
			t.invalidate(s)
		}
	}()
}

func (t *TokenInvalidationService) invalidate(s *types.Session) {
	t.logger.Warnf("invalidating session %s for user %s", s.SessionId, s.UserId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, tr := t.db.BeginTransaction(ctx)
	if tr != nil {
		t.logger.WithError(tr.SerializeError(errors.LvlDebug)).Error("failed to begin transaction")
		return
	}

	defer func() {
		tr = tx.Rollback(t.logger)
		if tr != nil {
			t.logger.WithError(tr.SerializeError(errors.LvlDebug)).Error("failed to rollback transaction")
		}
	}()

	tr = tx.Queries().DeleteSession(s.UserId, s.SessionId)
	if tr != nil {
		t.logger.WithError(tr.SerializeError(errors.LvlDebug)).Error("failed to delete session")
		return
	}

	tr = tx.Commit(t.logger)
	if tr != nil {
		t.logger.WithError(tr.SerializeError(errors.LvlDebug)).Error("failed to commit transaction")
		return
	}

	t.logger.Infof("successfully invalidated session %s for user %s", s.SessionId, s.UserId)
}

func (t *TokenInvalidationService) Channel() chan *types.Session {
	return t.receiveChannel
}
