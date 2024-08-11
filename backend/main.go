package main

import (
	"luna-backend/api"
	"luna-backend/sources"
	"luna-backend/sources/caldav"
	"luna-backend/util"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	caldavUrl, err := url.Parse(os.Getenv("CALDAV_URL"))
	if err != nil {
		panic(err)
	}

	util.Sources = []sources.Source{
		caldav.NewCaldavSource(
			caldavUrl,
			sources.BasicAuth(os.Getenv("CALDAV_USERNAME"), os.Getenv("CALDAV_PASSWORD")),
		),
	}

	api.Run()
}
