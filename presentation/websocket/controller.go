package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/services/client"
	"github.com/tiagorlampert/CHAOS/services/device"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type wsController struct {
	Configuration *environment.Configuration
	Logger        *logrus.Logger
	ClientService client.Service
	DeviceService device.Service
}

func NewController(
	configuration *environment.Configuration,
	logger *logrus.Logger,
	clientService client.Service,
	deviceService device.Service) {
	handler := &wsController{
		Configuration: configuration,
		Logger:        logger,
		ClientService: clientService,
		DeviceService: deviceService,
	}

	http.HandleFunc("/client", handler.client)
}

func (h *wsController) client(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("x-client")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Println("error connecting client:", err)
		return
	}

	err = h.ClientService.AddConnection(clientID, conn)
	if err != nil {
		h.Logger.Println("error adding client:", err)
		return
	}

	h.Logger.Println("Client connected: ", clientID)
}
