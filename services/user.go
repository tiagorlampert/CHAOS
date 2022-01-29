package services

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrInvalidPassword  = errors.New("invalid password")
)

type UpdateUserPasswordInput struct {
	Username, OldPassword, NewPassword string
}

type User interface {
	Insert(entities.User) error
	Login(username, password string) bool
	UpdatePassword(UpdateUserPasswordInput) error
}
