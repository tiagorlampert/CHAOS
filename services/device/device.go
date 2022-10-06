package device

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type Service interface {
	Insert(entities.Device) error
	FindAllConnected() ([]entities.Device, error)
	FindByMacAddress(address string) (*entities.Device, error)
}
