package util

import "luna-backend/sources"

// TODO: this will become better and will use a map once there is user logic and persistance
var Sources []sources.Source

func GetSource(id sources.SourceId) sources.Source {
	for _, source := range Sources {
		if source.GetId() == id {
			return source
		}
	}
	return nil
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
