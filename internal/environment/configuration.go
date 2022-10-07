package environment

import (
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Server    Server
	Database  Database
	SecretKey string `envconfig:"SECRET_KEY"`
}

type Server struct {
	Port string `envconfig:"PORT" validate:"required"`
}

type Database struct {
	Sqlite   Sqlite
	Postgres Postgres
}

func Load() (*Configuration, error) {
	configuration := &Configuration{}
	_ = readEnv(configuration)
	return configuration, configuration.Validate()
}

func readEnv(cfg interface{}) error {
	return envconfig.Process("", cfg)
}

func (c Configuration) Validate() error {
	return validator.New().Struct(c.Server)
}
