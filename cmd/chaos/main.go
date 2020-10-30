package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/handler/app"
	"github.com/tiagorlampert/CHAOS/internal/usecase/generate"
	"github.com/tiagorlampert/CHAOS/internal/util/ui"
	"github.com/tiagorlampert/CHAOS/pkg/system"
)

const AppName = "CHAOS"

var Version = "dev"

type App struct {
	handler handler.App
}

func main() {
	system.ValidateOS()
	ui.ShowMenu(Version)

	if err := NewApp().Run(); err != nil {
		log.WithField(`cause`, err).Fatal(fmt.Sprintf("error starting %s Application", AppName))
	}
}

func NewApp() *App {
	// Use Case
	generateUseCase := generate.NewGenerateUseCase()

	appHandler := app.NewAppHandler(generateUseCase)

	return &App{
		handler: appHandler,
	}
}

func (a *App) Run() error {
	a.handler.Handle()
	return nil
}
