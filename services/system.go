package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type System interface {
	Setup() (*entities.System, error)
	Load() (*entities.System, error)
	GetParams() (*entities.System, error)
	RefreshSecretKey() (string, error)
}
