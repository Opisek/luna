package main

import (
	"context"
	"fmt"
	"luna-backend/api"
	"luna-backend/common"
	"luna-backend/db"
	"luna-backend/interface/parsing"
	"luna-backend/log"
	"luna-backend/tasks"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var version string

func setupDirs(env *common.Environmental) error {
	err := os.MkdirAll(env.GetKeysPath(), 0660)
	return fmt.Errorf("could not create %v directory: %v", env.GetKeysPath(), err)
}

func setupConfig() (*logrus.Logger, *logrus.Entry, *common.CommonConfig, error) {
	var err error
	logger := log.NewLogger()
	mainLogger := logger.WithField("module", "main")

	env, err := common.ParseEnvironmental(mainLogger)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not parse environmental variables: %v", err)
	}

	commonConfig := &common.CommonConfig{
		Env: &env,
	}
	commonConfig.Version, err = common.ParseVersion(version)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("malformed binary version \"%v\": %v", version, err)
	}

	return logger, mainLogger, commonConfig, nil
}

func setupDb(commonConfig *common.CommonConfig, mainLogger *logrus.Entry, dbLogger *logrus.Entry) (*db.Database, error) {
	env := commonConfig.Env
	db := db.NewDatabase(env.DB_HOST, env.DB_PORT, env.DB_USERNAME, env.DB_PASSWORD, env.DB_DATABASE, commonConfig, parsing.GetPrimitivesParser(), dbLogger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	err = tx.Tables().InitializeVersionTable()
	if err != nil {
		return nil, fmt.Errorf("could not initialize version table: %v", err)
	}
	latestUsedVersion, err := tx.Queries().GetLatestVersion()
	if err != nil {
		tx.Rollback(mainLogger)
		return nil, err
	}
	if latestUsedVersion.IsGreaterThan(&commonConfig.Version) {
		err := fmt.Errorf("downgrades are not supported: database version %v is greater than binary version %v", latestUsedVersion.String(), commonConfig.Version.String())
		tx.Rollback(mainLogger)
		return nil, err
	}
	err = tx.Migrations().RunMigrations(&latestUsedVersion)
	if err != nil {
		err = fmt.Errorf("could not run migrations: %v", err)
		tx.Rollback(mainLogger)
		return nil, err
	}
	if !latestUsedVersion.IsEqualTo(&commonConfig.Version) {
		err = tx.Queries().UpdateVersion(commonConfig.Version)
		if err != nil {
			tx.Rollback(mainLogger)
			return nil, err
		}
	}
	return db, tx.Commit(mainLogger)
}

func createTask(name string, task func(*db.Transaction, *logrus.Entry) error, db *db.Database, cronLogger *logrus.Entry) func() {
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

		tx.Commit(cronLogger)
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
	var err error

	// Config
	logger, mainLogger, commonConfig, err := setupConfig()
	if err != nil {
		mainLogger.Errorf("could not set up config: %v", err)
		os.Exit(1)
	}

	// Directories
	setupDirs(commonConfig.Env)

	// Database
	dbLogger := logger.WithField("module", "database")
	dbReady := false
	var db *db.Database
	for i := 0; i < 5; i++ {
		db, err = setupDb(commonConfig, mainLogger, dbLogger)
		if err == nil {
			dbReady = true
			break
		}
		mainLogger.Warnf("could not set up database: %v", err)
		mainLogger.Warn("retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	if !dbReady {
		mainLogger.Errorf("could not set up database after 5 attempts: %v", err)
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

	// Wait for goroutines to finish
	var wg sync.WaitGroup
	startGoroutine(api.Run, &wg)
	startGoroutine(c.Start, &wg)
	wg.Wait()
}
