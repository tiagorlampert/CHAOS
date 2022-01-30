package environment

import (
	"github.com/go-playground/validator/v10"
	"os"
)

type Configuration struct {
	Server    Server
	Database  Database
	SecretKey string
}

type Server struct {
	Port string `validate:"required"`
}

type Database struct {
	Name string `validate:"required"`
}

func Load() *Configuration {
	return &Configuration{
		Server: Server{
			Port: os.Getenv(`PORT`),
		},
		Database: Database{
			Name: os.Getenv(`DATABASE_NAME`),
		},
		SecretKey: os.Getenv(`SECRET_KEY`),
	}
}

func (c Configuration) Validate() error {
	return validator.New().Struct(c)
}
