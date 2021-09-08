package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type Device interface {
	Insert(entities.Device) error
	GetAllAvailable() ([]entities.Device, error)
}
