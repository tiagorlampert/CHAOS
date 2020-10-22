package handler

import "github.com/tiagorlampert/CHAOS/internal/models"

type App interface {
	Handle()
}

type Server interface {
	HandleConnections()
	AcceptConnections()
	SetDevice(key string, con *models.Device)
	GetDevice(key string) (*models.Device, bool)
}

type Client interface {
	HandleConnection(hostname string, user string)
}
