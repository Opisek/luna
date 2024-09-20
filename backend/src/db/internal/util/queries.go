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

// parameters: we need a generic interface and then a method to get types.Id from that generic interface
func JoinIds[T any](idProviders []T, getId func(T) types.ID) []any {
	ids := make([]any, len(idProviders))

	for i, idProvider := range idProviders {
		ids[i] = getId(idProvider)
	}

	return ids
}
