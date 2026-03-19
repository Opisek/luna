package config

import (
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environmental struct {
	PUBLIC_URL url.URL `env:"PUBLIC_URL,required"`

	DB_URL      string `env:"DB_URL"`
	DB_HOST     string `env:"DB_HOST"`
	DB_PORT     uint16 `env:"DB_PORT" envDefault:"5432"`
	DB_USERNAME string `env:"DB_USERNAME"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_DATABASE string `env:"DB_DATABASE"`

	DATA_PATH string `env:"DATA_PATH" envDefault:"/data"`
	API_PORT  uint16 `env:"API_PORT" envDefault:"3000"`

	REQUEST_TIMEOUT_DEFAULT        time.Duration `env:"REQUEST_TIMEOUT_DEFAULT" envDefault:"15s"`
	REQUEST_TIMEOUT_AUTHENTICATION time.Duration `env:"REQUEST_TIMEOUT_AUTHENTICATION" envDefault:"15s"`

	DEVELOPMENT bool `env:"DEVELOPMENT" envDefault:"false"`
}

func ParseEnvironmental(logger *logrus.Entry) (Environmental, error) {

	if err := godotenv.Load(); err != nil {
		logger.Infof("could not load .env file: %v", err)
	}

	var environmental Environmental

	if err := env.Parse(&environmental); err != nil {
		return environmental, err
	}

	if err := environmental.Validate(); err != nil {
		return environmental, err
	}

	return environmental, nil
}

func (env *Environmental) Validate() error {
	if env.DB_URL == "" {
		if env.DB_HOST == "" {
			return fmt.Errorf("DB_HOST is required if DB_URL is not set")
		}
		if env.DB_USERNAME == "" {
			return fmt.Errorf("DB_USERNAME is required if DB_URL is not set")
		}
		if env.DB_PASSWORD == "" {
			return fmt.Errorf("DB_PASSWORD is required if DB_URL is not set")
		}
		if env.DB_DATABASE == "" {
			return fmt.Errorf("DB_DATABASE is required if DB_URL is not set")
		}
	}

	return nil
}

func (env *Environmental) getBasePath() string {
	return env.DATA_PATH
}

func (env *Environmental) GetKeysPath() string {
	return path.Join(env.getBasePath(), "keys")
}
