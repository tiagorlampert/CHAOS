package auth

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type Service interface {
	GetSecret() (string, error)
	GetAuthConfig() (*entities.Auth, error)
	RefreshSecret() (string, error)
}
