package repositories

import "github.com/tiagorlampert/CHAOS/entities"

type Auth interface {
	Insert(auth entities.Auth) error
	Update(auth entities.Auth) error
	First() (entities.Auth, error)
}
