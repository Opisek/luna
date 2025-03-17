package tasks

import (
	"encoding/json"
	"luna-backend/auth"
	"luna-backend/db"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/interface/protocols/ical"
	"sync"

	"github.com/sirupsen/logrus"
)

// This method is responsible for periodically fetching remote files to keep the
// local cache up to date, should the remote file be inaccessible when the user
// requests it later.
func RefetchIcalFiles(tx *db.Transaction, logger *logrus.Entry) *errors.ErrorTrace {
	settings, tr := tx.Queries().GetSourceSettingsByType("ical")

	if tr != nil {
		return tr.
			Append(errors.LvlDebug, "Could not refetch iCal files")
	}

	wg := sync.WaitGroup{}

	for _, setting := range settings {
		wg.Add(1)
		go func(setting []byte) {
			defer wg.Done()

			icalSourceSettings := &ical.IcalSourceSettings{}

			err := json.Unmarshal(setting, icalSourceSettings)
			if err != nil {
				logger.Errorf("could not unmarshal iCal settings: %v", err)
			}

			if icalSourceSettings.Location != "remote" {
				return
			}

			// We assume no authentication is needed for this file.
			// This will fail for users whose remote iCal files require authentication.
			// This will not be fixed in this task, because we don't want to expose users' encryption keys unnecessarily.
			// Instead, refetching of access-controlled iCal files might become an opt-in feature later on.
			file := files.NewRemoteFile(icalSourceSettings.Url, auth.NewNoAuth())
			tr = file.ForceFetchFromRemote(tx.Queries())

			if tr != nil {
				logger.Errorf("could not refetch iCal file %v: %v", icalSourceSettings.Url, tr.Serialize(errors.LvlDebug))
			}
		}(setting)
	}

	wg.Wait()
	return nil
}
