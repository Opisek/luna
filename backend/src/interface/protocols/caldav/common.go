package caldav

import "strings"

func unespaceString(s string) string {
	s = strings.ReplaceAll(s, "\\,", ",")
	s = strings.ReplaceAll(s, "\\:", ",")
	s = strings.ReplaceAll(s, "\\;", ";")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	return s
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

var ProductId string = "-//opisek.net//Luna//EN"
