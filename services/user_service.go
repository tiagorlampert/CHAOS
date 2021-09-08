package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
	repo "github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/shared/utils"
)

type userService struct {
	repository repo.User
}

func NewUser(repository repo.User) User {
	return &userService{
		repository: repository,
	}
}

func (u userService) Login(username, password string) bool {
	user, err := u.repository.Get(username)
	if err != nil {
		return false
	}
	return utils.PasswordsMatch(user.Password, password)
}

func (u userService) Create(input entities.User) error {
	_, err := u.repository.Get(input.Username)
	if err != nil {
		if err == repo.ErrNotFound {
			err := u.repository.Insert(input)
			if err != nil {
				return nil
			}
		}
		return err
	}
	return ErrUserAlreadyExist
}

func (u userService) UpdatePassword(input UpdateUserPasswordInput) error {
	user, err := u.repository.Get(input.Username)
	if err != nil {
		return err
	}

	if !utils.PasswordsMatch(user.Password, input.OldPassword) {
		return ErrInvalidPassword
	}

	hashedPassword, err := utils.HashAndSalt(input.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := u.repository.Update(user); err != nil {
		return err
	}
	return nil
}
