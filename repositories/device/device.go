package device

import (
	"github.com/tiagorlampert/CHAOS/entities"
	"time"
)

type Repository interface {
	Insert(device entities.Device) error
	Update(device entities.Device) error
	FindByMacAddress(address string) (*entities.Device, error)
	FindAll(updatedAt time.Time) ([]entities.Device, error)
}
