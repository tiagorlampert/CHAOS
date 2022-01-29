package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/shared/utilities"
)

const defaultPassword = "admin"

type authService struct {
	authRepository repositories.Auth
	userRepository repositories.User
}

func NewAuth(
	authRepository repositories.Auth,
	userRepository repositories.User) Auth {
	return &authService{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (s authService) Setup() (*entities.Auth, error) {
	auth, err := s.authRepository.First()
	if err == nil {
		return auth, nil
	}
	if err := s.authRepository.Insert(
		entities.Auth{SecretKey: utilities.GenerateRandomString(secretKeySize)},
	); err != nil {
		return nil, err
	}

	passwordHash, err := utilities.HashAndSalt(defaultPassword)
	if err != nil {
		return nil, err
	}
	if err := s.userRepository.Insert(entities.User{
		Username: "admin",
		Password: passwordHash,
	}); err != nil {
		return nil, err
	}
	return s.authRepository.First()
}

func (s authService) First() (*entities.Auth, error) {
	return s.authRepository.First()
}

func (s authService) RefreshSecret() (string, error) {
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
