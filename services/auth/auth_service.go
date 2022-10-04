package auth

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utils"
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

func (s authService) Setup() (*entities.Auth, error) {
	entry, err := s.AuthRepository.GetFirst()
	switch err {
	case nil, repositories.ErrNotFound:
		break
	default:
		return nil, err
	}

	hasProvidedSecretKey := len(strings.TrimSpace(s.SecretKey)) > 0
	if hasProvidedSecretKey {
		defer s.Logger.WithFields(logrus.Fields{"key": s.SecretKey}).
			Info("Using a provided secret key from environment variable")
	}

	if errors.Is(err, repositories.ErrNotFound) {
		authEntry := entities.Auth{}
		if hasProvidedSecretKey {
			authEntry.SecretKey = s.SecretKey
		} else {
			authEntry.SecretKey = utils.GenerateRandomString(secretKeySize)
		}

		if err := s.AuthRepository.Insert(authEntry); err != nil {
			return nil, err
		}
		return &authEntry, nil
	}

	if hasProvidedSecretKey && entry.SecretKey != s.SecretKey {
		entry.SecretKey = s.SecretKey

		if err := s.AuthRepository.Update(entry); err != nil {
			return nil, err
		}
	}
	return entry, nil
}

func (s authService) GetAuthConfig() (*entities.Auth, error) {
	return s.AuthRepository.GetFirst()
}

func (s authService) RefreshSecret() (string, error) {
	if len(s.SecretKey) != 0 {
		return "", fmt.Errorf("%s", ErrFailedRefreshProvidedSecretKey)
	}

	auth, err := s.AuthRepository.GetFirst()
	if err != nil {
		return "", err
	}
	if err := s.AuthRepository.Update(&entities.Auth{
		DBModel:   auth.DBModel,
		SecretKey: utils.GenerateRandomString(secretKeySize),
	}); err != nil {
		return "", err
	}
	auth, err = s.AuthRepository.GetFirst()
	if err != nil {
		return "", err
	}
	return auth.SecretKey, nil
}
