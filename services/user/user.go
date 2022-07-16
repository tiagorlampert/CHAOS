package user

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrInvalidPassword  = errors.New("invalid password")
)

type UpdateUserPasswordInput struct {
	Username    string
	OldPassword string
	NewPassword string
}

type Service interface {
	Insert(entities.User) error
	Login(username, password string) bool
	UpdatePassword(UpdateUserPasswordInput) error
	CreateDefaultUser() error
}
