package handler

import (
	"github.com/tiagorlampert/CHAOS/pkg/models"
)

type App interface {
	Handle() error
}

type Server interface {
	HandleConnections()
	AcceptConnections()
	SetDevice(key string, con *models.Device)
	GetDevice(key string) (*models.Device, bool)
}

type Client interface {
	HandleConnection() error
}
