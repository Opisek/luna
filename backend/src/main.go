package main

import (
	"context"
	"fmt"
	"luna-backend/api"
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/errors"
	"luna-backend/log"
	"luna-backend/parsing"
	"luna-backend/tasks"
	"luna-backend/types"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var version string

func setupDirs(env *config.Environmental) error {
	err := os.MkdirAll(env.GetKeysPath(), 0660)
	return fmt.Errorf("could not create %v directory: %v", env.GetKeysPath(), err)
}

func setupConfig() (*logrus.Logger, *logrus.Entry, *config.CommonConfig, *errors.ErrorTrace) {
	var err error
	logger := log.NewLogger()
	mainLogger := logger.WithField("module", "main")

	env, err := config.ParseEnvironmental(mainLogger)
	if err != nil {
		return logger, mainLogger, nil, errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse environmental variables")
	}

	commonConfig := &config.CommonConfig{
		Env: &env,
	}
	commonConfig.Version, err = types.ParseVersion(version)
	if err != nil {
		return logger, mainLogger, nil, errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse binary version %v", version)
	}

	return logger, mainLogger, commonConfig, nil
}

func setupDb(commonConfig *config.CommonConfig, mainLogger *logrus.Entry, dbLogger *logrus.Entry) (*db.Database, *errors.ErrorTrace) {
	env := commonConfig.Env
	db := db.NewDatabase(env.DB_HOST, env.DB_PORT, env.DB_USERNAME, env.DB_PASSWORD, env.DB_DATABASE, commonConfig, parsing.GetPrimitivesParser(), dbLogger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, tr := db.BeginTransaction(ctx)
	if tr != nil {
		return nil, tr
	}
	defer func() {
		rollbackErr := tx.Rollback(mainLogger)
		if rollbackErr != nil {
			mainLogger.Error(rollbackErr.Serialize(errors.LvlDebug))
		}
	}()

	// Verify version integrity
	err := tx.Tables().InitializeVersionTable()
	if err != nil {
		return nil, errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not initialize version table")
	}
	latestUsedVersion, tr := tx.Queries().GetLatestVersion()
	if tr != nil {
		return nil, tr
	}
	if latestUsedVersion.IsGreaterThan(&commonConfig.Version) {
		tr := errors.New().
			Append(errors.LvlDebug, "Database version %v is greater than binary version %v", latestUsedVersion.String(), commonConfig.Version.String()).
			Append(errors.LvlDebug, "Downgrades are not supported")
		return nil, tr
	}

	// Run migrations
	tr = tx.Migrations().RunMigrations(&latestUsedVersion)
	if tr != nil {
		return nil, tr
	}
	if !latestUsedVersion.IsEqualTo(&commonConfig.Version) {
		tr = tx.Queries().UpdateVersion(commonConfig.Version)
		if tr != nil {
			return nil, tr
		}
	}

	// Load global settings
	commonConfig.Settings, tr = tx.Queries().GetGlobalSettings()
	if tr != nil {
		return nil, tr
	}

	return db, tx.Commit(mainLogger)
}

func createTask(name string, task func(*db.Transaction, *logrus.Entry) *errors.ErrorTrace, db *db.Database, cronLogger *logrus.Entry) func() {
	return func() {
		cronLogger.Infof("running cron task %v", name)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tx, err := db.BeginTransaction(ctx)
		defer tx.Rollback(cronLogger)
		if err != nil {
			cronLogger.Errorf("failure creating database transaction for cron task %v: %v", name, err)
			return
		}

		err = task(tx, cronLogger.WithField("task", name))
		if err != nil {
			cronLogger.Errorf("failure running cron task %v: %v", name, err)
			return
		}

		err = tx.Commit(cronLogger)
		if err != nil {
			cronLogger.Errorf("failure committing transaction for cron task %v: %v", name, err)
			return
		}

		cronLogger.Infof("successfully finished cron task %v", name)
	}
}

func startGoroutine(f func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}

func main() {
	// Config
	logger, mainLogger, commonConfig, err := setupConfig()
	if err != nil {
		mainLogger.Errorf("could not set up config: %v", err.Serialize(errors.LvlDebug))
		os.Exit(1)
	}

	// Directories
	setupDirs(commonConfig.Env)

	// Database
	dbLogger := logger.WithField("module", "database")
	dbReady := false
	var db *db.Database
	for range 5 {
		db, err = setupDb(commonConfig, mainLogger, dbLogger)
		if err == nil {
			dbReady = true
			break
		}
		mainLogger.Warnf("could not set up database: %v", err.Serialize(errors.LvlDebug))
		mainLogger.Warn("retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	if !dbReady {
		mainLogger.Errorf("could not set up database after 5 attempts: %v", err.Serialize(errors.LvlDebug))
		os.Exit(1)
	}

	// Api Server
	apiLogger := logger.WithField("module", "api")
	api := api.NewApi(db, commonConfig, apiLogger)
	mainLogger.Infof("started luna-backend %s", commonConfig.Version.String())

	// Scheduled tasks
	cronLogger := logger.WithField("module", "cron")
	c := cron.New()
	c.AddFunc("*/30 * * * *", createTask("RefetchIcalFiles", tasks.RefetchIcalFiles, db, cronLogger))
	c.AddFunc("0 * * * *", createTask("DeleteExpiredShortLivedSessions", tasks.DeleteStaleShortLivedSessions, db, cronLogger))
	c.AddFunc("0 0 * * *", createTask("DeleteExpiredLongLivedSessions", tasks.DeleteStaleLongLivedSessions, db, cronLogger))

	// Wait for goroutines to finish
	var wg sync.WaitGroup
	startGoroutine(api.Run, &wg)
	startGoroutine(c.Start, &wg)
	wg.Wait()
}
