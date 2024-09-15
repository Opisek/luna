package caldav

import "strings"

func unespaceString(s string) string {
	s = strings.ReplaceAll(s, "\\,", ",")
	s = strings.ReplaceAll(s, "\\;", ";")
	s = strings.ReplaceAll(s, "\\n", "\n")
	return s
}
