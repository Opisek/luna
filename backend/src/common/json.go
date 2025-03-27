package common

import (
	"fmt"
	"strconv"
)

// Common JSON marshal and unmarshal functions

func MarshalBool(value bool) []byte {
	if value {
		return []byte("true")
	} else {
		return []byte("false")
	}
}

func UnmarshalBool(data []byte) (bool, error) {
	switch string(data) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("could not parse boolean: %v", string(data))
	}
}

func MarshalInt(value int) []byte {
	return []byte(strconv.Itoa(value))
}

func UnmarshalInt(data []byte) (int, error) {
	value, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("could not parse integer: %v", err)
	}
	return value, nil
}

func MarshalString(value string) []byte {
	return []byte("\"" + value + "\"")
}

func UnmarshalString(data []byte) (string, error) {
	// TODO: these checks are not exhausting: any unescaped " inside the string should result in a failure
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return "", fmt.Errorf("could not parse string: %v", string(data))
	}
	return string(data[1 : len(data)-1]), nil
}
