package device

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type Service interface {
	Insert(entities.Device) error
	FindAll() ([]entities.Device, error)
}
