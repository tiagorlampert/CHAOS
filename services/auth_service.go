package services

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utilities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"strings"
)

type authService struct {
	logger         *logrus.Logger
	secretKey      string
	authRepository repositories.Auth
}

func NewAuth(
	logger *logrus.Logger,
	secretKey string,
	authRepository repositories.Auth) Auth {
	return &authService{
		logger:         logger,
		authRepository: authRepository,
		secretKey:      strings.TrimSpace(secretKey),
	}
}

func (s authService) Setup() (*entities.Auth, error) {
	auth, err := s.authRepository.First()
	switch err {
	case nil, repositories.ErrNotFound:
		break
	default:
		return nil, err
	}

	hasProvidedSecretKey := len(s.secretKey) > 0
	if hasProvidedSecretKey {
		defer s.logger.WithFields(logrus.Fields{"key": s.secretKey}).
			Info("Using a provided secret key from environment variable")
	}

	if errors.Is(err, repositories.ErrNotFound) {
		dummyAuth := entities.Auth{}
		if hasProvidedSecretKey {
			dummyAuth.SecretKey = s.secretKey
		} else {
			dummyAuth.SecretKey = utilities.GenerateRandomString(secretKeySize)
		}
		return &dummyAuth, s.authRepository.Insert(dummyAuth)
	}

	if hasProvidedSecretKey && auth.SecretKey != s.secretKey {
		auth.SecretKey = s.secretKey
		return &auth, s.authRepository.Update(auth)
	}
	return &auth, nil
}

func (s authService) First() (entities.Auth, error) {
	return s.authRepository.First()
}

func (s authService) RefreshSecret() (string, error) {
	if len(s.secretKey) != 0 {
		return "", fmt.Errorf("%s", ErrFailedRefreshProvidedSecretKey)
	}

	auth, err := s.authRepository.First()
	if err != nil {
		return "", err
	}
	if err := s.authRepository.Update(entities.Auth{
		DBModel:   auth.DBModel,
		SecretKey: utilities.GenerateRandomString(secretKeySize),
	}); err != nil {
		return "", err
	}
	auth, err = s.authRepository.First()
	if err != nil {
		return "", err
	}
	return auth.SecretKey, nil
}
