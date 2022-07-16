package auth

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
)

var ErrFailedRefreshProvidedSecretKey = errors.New("the secret key provided from environment variable cannot be redefined")

type Service interface {
	Setup() (*entities.Auth, error)
	GetAuthConfig() (entities.Auth, error)
	RefreshSecret() (string, error)
}
