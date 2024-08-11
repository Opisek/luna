package main

import (
	"luna-backend/api"
	"luna-backend/sources"
	"luna-backend/sources/caldav"
	"luna-backend/sources/ical"
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

	icalUrl, err := url.Parse(os.Getenv("ICAL_URL"))
	if err != nil {
		panic(err)
	}

	util.Sources = []sources.Source{
		caldav.NewCaldavSource(&caldav.CaldavSettings{
			Url:      caldavUrl,
			Username: os.Getenv("CALDAV_USERNAME"),
			Password: os.Getenv("CALDAV_PASSWORD"),
		}),
		ical.NewIcalSource(&ical.IcalSettings{
			Url: icalUrl,
		}),
	}

	api.Run()
}
