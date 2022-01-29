package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type Auth interface {
	Setup() (*entities.Auth, error)
	First() (*entities.Auth, error)
	RefreshSecret() (string, error)
}
