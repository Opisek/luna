package common

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environmental struct {
	DB_HOST     string
	DB_PORT     uint16
	DB_USERNAME string
	DB_PASSWORD string
	DB_DATABASE string
}

func ParseEnvironmental(logger *logrus.Entry) (Environmental, error) {
	env := Environmental{}

	err := godotenv.Load()
	if err != nil {
		logger.Warnf("could not load .env file: %v", err)
		//return env, errors.Join(errors.New("could not load .env file: "), err)
	}

	reflected := reflect.Indirect(reflect.ValueOf(&env))
	for i := 0; i < reflected.NumField(); i += 1 {
		field := reflected.Type().Field(i)
		fieldName := field.Name
		fieldType := field.Type.Name()

		fieldValueRaw := os.Getenv(fieldName)

		if fieldValueRaw == "" {
			err := fmt.Errorf("environmental variable %v is missing", fieldName)
			return env, err
		}

		switch fieldType {
		case "string":
			reflected.Field(i).SetString(fieldValueRaw)
		case "uint16":
			fieldValue, err := strconv.ParseUint(fieldValueRaw, 10, 16)
			if err != nil {
				err := fmt.Errorf("environmental variable %v is malformed", fieldName)
				return env, err
			}
			reflected.Field(i).SetUint(fieldValue)
		default:
			err := fmt.Errorf("unsupported type %v for environmental variable %v", fieldType, fieldName)
			return env, err
		}
	}

	return env, nil
}
