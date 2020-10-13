package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/handler"
	"github.com/tiagorlampert/CHAOS/client/app/handler/connection"
)

type App struct {
	Handler handler.Handler
}

func NewApp(address, port string) *App {
	connectionHandler := connection.NewConnectionHandler(address, port)

	return &App{
		Handler: connectionHandler,
	}
}

func (app *App) Run() error {
	if err := app.Handler.Handle(); err != nil {
		log.WithField("cause", err.Error()).Error("error handling app")
		return err
	}
	return nil
}
