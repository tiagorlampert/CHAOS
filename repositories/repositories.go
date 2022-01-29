package repositories

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
	"time"
)

var ErrNotFound = errors.New("not found")

type Auth interface {
	Insert(auth entities.Auth) error
	Update(auth entities.Auth) error
	First() (*entities.Auth, error)
}

type User interface {
	Insert(user entities.User) error
	Update(user *entities.User) error
	Get(username string) (*entities.User, error)
}

type Device interface {
	Insert(device entities.Device) error
	Update(device entities.Device) error
	GetByMacAddress(macAddr string) (*entities.Device, error)
	List(updatedAt time.Time) ([]entities.Device, error)
}
