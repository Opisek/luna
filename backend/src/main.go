package main

import (
	"luna-backend/api"
	"luna-backend/common"
	"luna-backend/db"
	"luna-backend/log"
	"os"
)

var version string

func main() {
	var err error

	//
	// Config
	//
	logger := log.NewLogger()
	mainLogger := logger.WithField("module", "main")

	commonConfig := &common.CommonConfig{}
	commonConfig.Version, err = common.ParseVersion(version)
	if err != nil {
		mainLogger.Errorf("malformed binary version \"%v\": %v", version, err)
		os.Exit(1)
	}

	env, err := common.ParseEnvironmental(mainLogger)
	if err != nil {
		mainLogger.Errorf("could not parse environmental variables: %v", err)
		os.Exit(1)
	}

	//
	// Database
	//
	dbLogger := logger.WithField("module", "database")
	db := db.NewDatabase(env.DB_HOST, env.DB_PORT, env.DB_USERNAME, env.DB_PASSWORD, env.DB_DATABASE, commonConfig, dbLogger)
	// Connect to the database
	err = db.Connect()
	if err != nil {
		mainLogger.Error(err)
		os.Exit(1)
	}
	// Run migrations
	latestUsedVersion, err := db.GetLatestVersion()
	if err != nil {
		mainLogger.Error(err)
		os.Exit(1)
	}
	if latestUsedVersion.IsGreaterThan(&commonConfig.Version) {
		mainLogger.Errorf("downgrades are not supported: database version %v is greater than binary version %v", latestUsedVersion.String(), commonConfig.Version.String())
		os.Exit(1)
	}
	err = db.RunMigrations(&latestUsedVersion)
	if err != nil {
		mainLogger.Errorf("could not run migrations: %v", err)
		os.Exit(1)
	}
	if !latestUsedVersion.IsEqualTo(&commonConfig.Version) {
		err = db.UpdateVersion(commonConfig.Version)
		if err != nil {
			mainLogger.Error(err)
			os.Exit(1)
		}
	}

	//caldavUrl, err := url.Parse(os.Getenv("CALDAV_URL"))
	//if err != nil {
	//	panic(err)
	//}

	//util.Sources = []sources.Source{
	//	caldav.NewCaldavSource(
	//		caldavUrl,
	//		auth.NewBasicAuth(os.Getenv("CALDAV_USERNAME"), os.Getenv("CALDAV_PASSWORD")),
	//	),
	//}

	// Api Server
	apiLogger := logger.WithField("module", "api")
	api := api.NewApi(db, commonConfig, apiLogger)
	mainLogger.Infof("started luna-backend %s", commonConfig.Version.String())
	api.Run()
}
