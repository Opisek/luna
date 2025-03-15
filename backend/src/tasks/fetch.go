package tasks

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/db"
	"luna-backend/files"
	"luna-backend/interface/protocols/ical"

	"github.com/sirupsen/logrus"
)

// This method is responsible for periodically fetching remote files to keep the
// local cache up to date, should the remote file be inaccessible when the user
// requests it later.
func RefetchIcalFiles(tx *db.Transaction, logger *logrus.Entry) error {
	settings, err := tx.Queries().GetSourceSettingsByType("ical")

	if err != nil {
		return fmt.Errorf("could not refetch iCal files: %v", err)
	}

	for _, setting := range settings {
		icalSourceSettings := &ical.IcalSourceSettings{}
		err = json.Unmarshal(setting, icalSourceSettings)
		if err != nil {
			return fmt.Errorf("could not unmarshal ical settings: %v", err)
		}

		if icalSourceSettings.Location != "remote" {
			continue
		}

		go func(sourceSettings *ical.IcalSourceSettings) {
			// We assume no authentication is needed for this file.
			// This will fail for users whose remote iCal files require authentication.
			// This will not be fixed in this task, because we don't want to expose users' encryption keys unnecessarily.
			// Instead, refetching of access-controlled iCal files might become an opt-in feature later on.
			file := files.NewRemoteFile(sourceSettings.Url, auth.NewNoAuth())
			err := file.ForceFetchFromRemote(tx.Queries())

			if err != nil {
				logger.Errorf("could not refetch iCal file %v: %v", sourceSettings.Url, err)
			}
		}(icalSourceSettings)
	}

	return nil
}
