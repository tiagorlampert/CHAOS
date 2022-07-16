package auth

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/repositories"
	"strings"
)

const secretKeySize = 50

type authService struct {
	Logger         *logrus.Logger
	SecretKey      string
	AuthRepository repositories.Auth
}

func NewAuthService(
	logger *logrus.Logger,
	secretKey string,
	authRepository repositories.Auth,
) Service {
	return &authService{
		Logger:         logger,
		AuthRepository: authRepository,
		SecretKey:      strings.TrimSpace(secretKey),
	}
}

func (s authService) Setup() (*entities.Auth, error) {
	auth, err := s.AuthRepository.First()
	switch err {
	case nil, repositories.ErrNotFound:
		break
	default:
		return nil, err
	}

	hasProvidedSecretKey := len(s.SecretKey) > 0
	if hasProvidedSecretKey {
		defer s.Logger.WithFields(logrus.Fields{"key": s.SecretKey}).
			Info("Using a provided secret key from environment variable")
	}

	if errors.Is(err, repositories.ErrNotFound) {
		dummyAuth := entities.Auth{}
		if hasProvidedSecretKey {
			dummyAuth.SecretKey = s.SecretKey
		} else {
			dummyAuth.SecretKey = utils.GenerateRandomString(secretKeySize)
		}
		return &dummyAuth, s.AuthRepository.Insert(dummyAuth)
	}

	if hasProvidedSecretKey && auth.SecretKey != s.SecretKey {
		auth.SecretKey = s.SecretKey
		return &auth, s.AuthRepository.Update(auth)
	}
	return &auth, nil
}

func (s authService) First() (entities.Auth, error) {
	return s.AuthRepository.First()
}

func (s authService) RefreshSecret() (string, error) {
	if len(s.SecretKey) != 0 {
		return "", fmt.Errorf("%s", ErrFailedRefreshProvidedSecretKey)
	}

	auth, err := s.AuthRepository.First()
	if err != nil {
		return "", err
	}
	if err := s.AuthRepository.Update(entities.Auth{
		DBModel:   auth.DBModel,
		SecretKey: utils.GenerateRandomString(secretKeySize),
	}); err != nil {
		return "", err
	}
	auth, err = s.AuthRepository.First()
	if err != nil {
		return "", err
	}
	return auth.SecretKey, nil
}
