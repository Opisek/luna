package util

import (
	"fmt"
	"luna-backend/types"
	"strings"
)

func GenerateArgList(start int, count int) string {
	args := make([]string, count)

	for i := 0; i < count; i++ {
		args[i] = fmt.Sprintf("$%d", start+i)
	}

	return strings.Join(args, ", ")
}

func GenerateSetList(start int, columns []string) string {
	set := make([]string, len(columns))

	for i, column := range columns {
		set[i] = fmt.Sprintf("%s = $%d", column, start+i)
	}

	return strings.Join(set, ", ")
}

// parameters: we need a generic interface and then a method to get types.Id from that generic interface
func JoinIds[T any](idProviders []T, getId func(T) types.ID) []any {
	ids := make([]any, len(idProviders))

	for i, idProvider := range idProviders {
		ids[i] = getId(idProvider)
	}

	return ids
}
