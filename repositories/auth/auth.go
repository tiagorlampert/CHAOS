package auth

import "github.com/tiagorlampert/CHAOS/entities"

type Repository interface {
	Insert(auth entities.Auth) error
	Update(auth *entities.Auth) error
	GetFirst() (*entities.Auth, error)
}
