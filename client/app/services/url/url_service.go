package url

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"runtime"
	"strings"
)

type URLService struct {
	Terminal services.Terminal
}

func NewURLService(terminalService services.Terminal) services.URL {
	return &URLService{terminalService}
}

func (u URLService) OpenURL(url string) error {
	var cmdOut string
	switch runtime.GOOS {
	case `windows`:
		cmdOut = u.Terminal.Run(fmt.Sprintf("start %s", url), 10)
	case `linux`:
		cmdOut = u.Terminal.Run(fmt.Sprintf("xdg-open %s", url), 10)
	default:
		return services.ErrUnsupportedPlatform
	}
	if strings.Contains(strings.ToLower(cmdOut), "failed") {
		return fmt.Errorf("%s", cmdOut)
	}
	return nil
}
