package auth

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utils/random"
	"github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/repositories/auth"
	"strings"
)

const secretKeySize = 50

type authService struct {
	Logger         *logrus.Logger
	SecretKey      string
	AuthRepository auth.Repository
}

func NewAuthService(
	logger *logrus.Logger,
	secretKey string,
	authRepository auth.Repository,
) Service {
	return &authService{
		Logger:         logger,
		AuthRepository: authRepository,
		SecretKey:      strings.TrimSpace(secretKey),
	}
}

func (a authService) GetSecret() (string, error) {
	setup, err := a.Setup()
	if err != nil {
		return "", err
	}
	return setup.SecretKey, nil
}

func (a authService) Setup() (*entities.Auth, error) {
	auth, err := a.AuthRepository.GetFirst()
	switch err {
	case nil, repositories.ErrNotFound:
		break
	default:
		return nil, err
	}

	hasProvidedSecretKey := len(strings.TrimSpace(a.SecretKey)) > 0
	if hasProvidedSecretKey {
		defer a.Logger.WithFields(logrus.Fields{"key": a.SecretKey}).
			Info("Using a provided secret key from environment variable")
	}

	if errors.Is(err, repositories.ErrNotFound) {
		authEntry := entities.Auth{}
		if hasProvidedSecretKey {
			authEntry.SecretKey = a.SecretKey
		} else {
			authEntry.SecretKey = random.GenerateString(secretKeySize)
		}

		if err := a.AuthRepository.Insert(authEntry); err != nil {
			return nil, err
		}
		return &authEntry, nil
	}

	if hasProvidedSecretKey && auth.SecretKey != a.SecretKey {
		auth.SecretKey = a.SecretKey

		if err := a.AuthRepository.Update(auth); err != nil {
			return nil, err
		}
	}
	return auth, nil
}

func (a authService) GetAuthConfig() (*entities.Auth, error) {
	return a.AuthRepository.GetFirst()
}

func (a authService) RefreshSecret() (string, error) {
	auth, err := a.AuthRepository.GetFirst()
	if err != nil {
		return "", err
	}
	auth.SecretKey = random.GenerateString(secretKeySize)
	err = a.AuthRepository.Update(auth)
	return auth.SecretKey, err
}
