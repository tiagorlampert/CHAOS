package environment

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" validate:"required"`
	Port     string `envconfig:"POSTGRES_PORT" validate:"required"`
	User     string `envconfig:"POSTGRES_USER" validate:"required"`
	Password string `envconfig:"POSTGRES_PASSWORD" validate:"required"`
	Database string `envconfig:"POSTGRES_DATABASE" validate:"required"`
	SSLMode  string `envconfig:"POSTGRES_SSL_MODE"`
}

func (p Postgres) BuildConnectionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=%s",
		p.User, p.Password, p.Host, p.Database, p.Port, p.SSLMode)
}

func (p Postgres) IsValid() bool {
	return validator.New().Struct(p) == nil
}
