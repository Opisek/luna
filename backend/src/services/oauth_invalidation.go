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

type OauthInvalidationService struct {
	receiveChannel chan types.ID
	db             *db.Database
	commonConfig   *config.CommonConfig
	logger         *logrus.Entry
}

func NewOauthInvalidationService(db *db.Database, commonConfig *config.CommonConfig, logger *logrus.Entry) *OauthInvalidationService {
	service := OauthInvalidationService{
		receiveChannel: make(chan types.ID),
		db:             db,
		commonConfig:   commonConfig,
		logger:         logger,
	}

	commonConfig.OauthInvalidationChannel = service.Channel()

	return &service
}

func (t *OauthInvalidationService) Start() {
	go func() {
		for s := range t.receiveChannel {
			t.invalidate(s)
		}
	}()
}

func (t *OauthInvalidationService) invalidate(tokenId types.ID) {
	t.logger.Warnf("invalidating OAuth 2.0 tokens %v", tokenId)

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

	tr = tx.Queries().DeleteOauthTokens(tokenId)
	if tr != nil {
		t.logger.WithError(tr.SerializeError(errors.LvlDebug)).Error("failed to delete session")
		return
	}

	tr = tx.Commit(t.logger)
	if tr != nil {
		t.logger.WithError(tr.SerializeError(errors.LvlDebug)).Error("failed to commit transaction")
		return
	}

	t.logger.Infof("successfully invalidated OAuth 2.0 tokens %v", tokenId)
}

func (t *OauthInvalidationService) Channel() chan types.ID {
	return t.receiveChannel
}
