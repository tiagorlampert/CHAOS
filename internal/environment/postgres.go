package environment

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Postgres struct {
	Host     string `yaml:"host" json:"host" envconfig:"POSTGRES_HOST" validate:"required"`
	Port     string `yaml:"port" json:"port" envconfig:"POSTGRES_PORT" validate:"required"`
	User     string `yaml:"user" json:"user" envconfig:"POSTGRES_USER" validate:"required"`
	Password string `yaml:"password" json:"password" envconfig:"POSTGRES_PASSWORD" validate:"required"`
	Database string `yaml:"database" json:"database" envconfig:"POSTGRES_DATABASE" validate:"required"`
	SSLMode  string `yaml:"ssl-mode" json:"ssl_mode" envconfig:"POSTGRES_SSL_MODE" validate:"required"`
}

func (p Postgres) BuildConnectionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=%s",
		p.User, p.Password, p.Host, p.Database, p.Port, p.SSLMode)
}

func (p Postgres) IsValid() bool {
	return validator.New().Struct(p) == nil
}
