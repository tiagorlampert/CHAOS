package repositories

import "github.com/tiagorlampert/CHAOS/entities"

type User interface {
	Insert(user entities.User) error
	Update(user *entities.User) error
	FindByUsername(username string) (*entities.User, error)
}
