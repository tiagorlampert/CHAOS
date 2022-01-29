package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
	repo "github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/shared/utilities"
)

type userService struct {
	repository repo.User
}

func NewUser(repository repo.User) User {
	return &userService{repository: repository}
}

func (u userService) Login(username, password string) bool {
	user, err := u.repository.Get(username)
	if err != nil {
		return false
	}
	return utilities.PasswordsMatch(user.Password, password)
}

func (u userService) Insert(input entities.User) error {
	_, err := u.repository.Get(input.Username)
	switch err {
	case repo.ErrNotFound:
		return u.repository.Insert(input)
	default:
		return ErrUserAlreadyExist
	}
}

func (u userService) UpdatePassword(input UpdateUserPasswordInput) error {
	user, err := u.repository.Get(input.Username)
	if err != nil {
		return err
	}
	if !utilities.PasswordsMatch(user.Password, input.OldPassword) {
		return ErrInvalidPassword
	}

	passwordHash, err := utilities.HashAndSalt(input.NewPassword)
	if err != nil {
		return err
	}
	user.Password = passwordHash
	return u.repository.Update(user)
}
