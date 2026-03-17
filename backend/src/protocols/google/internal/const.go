package google

import "luna-backend/types"

func ApiUrl() *types.Url {
	return types.NewUrlSafe("https://www.googleapis.com/calendar/v3")
}
