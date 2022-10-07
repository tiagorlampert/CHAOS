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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientID := r.Header.Get("x-client")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Println("error connecting client:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ClientService.AddConnection(clientID, ws)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Println("error adding client:", err)
		return
	}

	h.Logger.Println("Client connected: ", clientID)
}
