package util

import "strings"

var illegalSubstrings = []string{
	"\x00",
	"'",
	"\"",
	"\b",
	"\n",
	"\r",
	"\t",
	"\\",
	"%",
	"_",
	"DROP",
	"DELETE",
	"INSERT",
	"UPDATE",
	"SELECT",
	";",
}

func isSafe(str string) bool {
	for _, subStr := range illegalSubstrings {
		if strings.Contains(str, subStr) {
			return false
		}
	}
	return true
}
