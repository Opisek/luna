package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environmental struct {
	PUBLIC_URL string

	DB_HOST     string
	DB_PORT     uint16
	DB_USERNAME string
	DB_PASSWORD string
	DB_DATABASE string

	DATA_PATH string
	API_PORT  uint16

	REQUEST_TIMEOUT_DEFAULT        time.Duration
	REQUEST_TIMEOUT_AUTHENTICATION time.Duration

	DEVELOPMENT bool
}

var defaultEnv = Environmental{
	DATA_PATH:                      "/data",
	API_PORT:                       3000,
	REQUEST_TIMEOUT_DEFAULT:        5 * time.Second,
	REQUEST_TIMEOUT_AUTHENTICATION: 5 * time.Second,
	DEVELOPMENT:                    false,
}

func ParseEnvironmental(logger *logrus.Entry) (Environmental, error) {
	env := defaultEnv

	err := godotenv.Load()
	if err != nil {
		logger.Infof("could not load .env file: %v", err)
	}

	reflected := reflect.Indirect(reflect.ValueOf(&env))
	for i := 0; i < reflected.NumField(); i += 1 {
		field := reflected.Type().Field(i)
		fieldName := field.Name
		fieldType := field.Type.Name()

		fieldValueRaw := os.Getenv(fieldName)

		if fieldValueRaw == "" {
			if reflected.Field(i).IsZero() && fieldType != "bool" {
				err := fmt.Errorf("environmental variable %v is missing", fieldName)
				return env, err
			} else {
				logger.Warnf("environmental variable %v is missing, using default value %v", fieldName, reflected.Field(i).Interface())
				continue
			}
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
		case "Duration":
			fieldValue, err := strconv.ParseUint(fieldValueRaw, 10, 32)
			if err != nil {
				err := fmt.Errorf("environmental variable %v is malformed", fieldName)
				return env, err
			}
			duration := time.Duration(fieldValue * uint64(time.Second))
			reflected.Field(i).SetInt(duration.Nanoseconds())
		case "bool":
			fieldValue, err := strconv.ParseBool(fieldValueRaw)
			if err != nil {
				err := fmt.Errorf("environmental variable %v is malformed", fieldName)
				return env, err
			}
			reflected.Field(i).SetBool(fieldValue)
		default:
			err := fmt.Errorf("unsupported type %v for environmental variable %v", fieldType, fieldName)
			return env, err
		}
	}

	return env, nil
}

func (env *Environmental) getBasePath() string {
	return env.DATA_PATH
}

func (env *Environmental) GetKeysPath() string {
	return env.getBasePath() + "/keys"
}
