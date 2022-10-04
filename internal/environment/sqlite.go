package environment

import "github.com/go-playground/validator/v10"

type Sqlite struct {
	DatabaseName string `envconfig:"DATABASE_NAME" validate:"required"`
}

func (s Sqlite) IsValid() bool {
	return validator.New().Struct(s) == nil
}
