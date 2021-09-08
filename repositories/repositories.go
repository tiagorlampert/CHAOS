package repositories

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
	"time"
)

var ErrNotFound = errors.New("not found")

type System interface {
	Insert(system entities.System) error
	Update(system entities.System) (*entities.System, error)
	Get() (*entities.System, error)
}

type User interface {
	Insert(user entities.User) error
	Update(user *entities.User) error
	Get(username string) (*entities.User, error)
}

type Device interface {
	Insert(device entities.Device) error
	Update(device entities.Device) error
	Get(macAddress string) (*entities.Device, error)
	List(dateTime time.Time) ([]entities.Device, error)
}
