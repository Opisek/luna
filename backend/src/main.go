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

func setupDirs() error {
	os.MkdirAll("/data/keys", 0660)
	return nil
}

func setupConfig() (*logrus.Logger, *logrus.Entry, *common.CommonConfig, *common.Environmental, error) {
	var err error
	logger := log.NewLogger()
	mainLogger := logger.WithField("module", "main")

	commonConfig := &common.CommonConfig{}
	commonConfig.Version, err = common.ParseVersion(version)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("malformed binary version \"%v\": %v", version, err)
	}

	env, err := common.ParseEnvironmental(mainLogger)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("could not parse environmental variables: %v", err)
	}

	return logger, mainLogger, commonConfig, &env, nil
}

func setupDb(env *common.Environmental, commonConfig *common.CommonConfig, mainLogger *logrus.Entry, dbLogger *logrus.Entry) (*db.Database, error) {
	db := db.NewDatabase(env.DB_HOST, env.DB_PORT, env.DB_USERNAME, env.DB_PASSWORD, env.DB_DATABASE, commonConfig, dbLogger)

	tx, err := db.BeginTransaction()
	if err != nil {
		return nil, err
	}
	latestUsedVersion, err := tx.GetLatestVersion()
	if err != nil {
		tx.Rollback(mainLogger)
		return nil, err
	}
	if latestUsedVersion.IsGreaterThan(&commonConfig.Version) {
		err := fmt.Errorf("downgrades are not supported: database version %v is greater than binary version %v", latestUsedVersion.String(), commonConfig.Version.String())
		tx.Rollback(mainLogger)
		return nil, err
	}
	err = tx.RunMigrations(&latestUsedVersion)
	if err != nil {
		err = fmt.Errorf("could not run migrations: %v", err)
		tx.Rollback(mainLogger)
		return nil, err
	}
	if !latestUsedVersion.IsEqualTo(&commonConfig.Version) {
		err = tx.UpdateVersion(commonConfig.Version)
		if err != nil {
			tx.Rollback(mainLogger)
			return nil, err
		}
	}
	return db, tx.Commit(mainLogger)
}

func main() {
	var err error

	// Directories
	setupDirs()

	// Config
	logger, mainLogger, commonConfig, env, err := setupConfig()
	if err != nil {
		mainLogger.Errorf("could not set up config: %v", err)
		os.Exit(1)
	}

	// Database
	dbLogger := logger.WithField("module", "database")
	dbReady := false
	var db *db.Database
	for i := 0; i < 5; i++ {
		db, err = setupDb(env, commonConfig, mainLogger, dbLogger)
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
