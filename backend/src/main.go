package main

import (
	"fmt"
	"luna-backend/api"
	"luna-backend/common"
	"luna-backend/db"
	"luna-backend/log"
	"os"
	"time"

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
	db := db.NewDatabase(env.DB_HOST, env.DB_PORT, env.DB_USERNAME, env.DB_PASSWORD, env.DB_DATABASE, commonConfig, dbLogger)

	tx, err := db.BeginTransaction()
	if err != nil {
		return nil, err
	}
	err = tx.Tables().InitalizeVersionTable()
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
	api.Run()
}
