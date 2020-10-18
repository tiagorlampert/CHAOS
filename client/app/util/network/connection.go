package network

import (
	"bufio"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

func NewConnection(address, port string) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		log.WithField("cause", err.Error()).Fatal("error creating new connection")
	}
	return conn
}

func Send(conn net.Conn, input []byte) error {
	encoded := base64.StdEncoding.EncodeToString(input)
	if err := Write(conn, encoded); err != nil {
		log.WithField("cause", err.Error()).Error("error sending command to client")
		return err
	}
	return nil
}

func Read(conn net.Conn) ([]byte, error) {
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
	_, err := conn.Write([]byte(v + "\n"))
	return err
}
