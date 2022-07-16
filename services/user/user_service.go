package user

import (
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	repo "github.com/tiagorlampert/CHAOS/repositories"
)

const (
	defaultUser     = "admin"
	defaultPassword = "admin"
)

type userService struct {
	repository repo.User
}

func NewUserService(repository repo.User) Service {
	return &userService{repository: repository}
}

func (u userService) Login(username, password string) bool {
	user, err := u.repository.FindByUsername(username)
	if err != nil {
		return false
	}
	return utils.PasswordsMatch(user.Password, password)
}

func (u userService) Insert(input entities.User) error {
	_, err := u.repository.FindByUsername(input.Username)
	switch err {
	case repo.ErrNotFound:
		return u.repository.Insert(input)
	default:
		return ErrUserAlreadyExist
	}
}

func (u userService) UpdatePassword(input UpdateUserPasswordInput) error {
	user, err := u.repository.FindByUsername(input.Username)
	if err != nil {
		return err
	}
	if !utils.PasswordsMatch(user.Password, input.OldPassword) {
		return ErrInvalidPassword
	}

	passwordHash, err := utils.HashAndSalt(input.NewPassword)
	if err != nil {
		return err
	}
	user.Password = passwordHash
	return u.repository.Update(user)
}

func (u userService) CreateDefaultUser() error {
	_, err := u.repository.FindByUsername(defaultUser)
	switch err {
	case repo.ErrNotFound:
		break
	default:
		return err
	}

	passwordHash, err := utils.HashAndSalt(defaultPassword)
	if err != nil {
		return err
	}
	return u.repository.Insert(entities.User{
		Username: defaultUser,
		Password: passwordHash,
	})
}
