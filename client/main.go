package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app"
)

var (
	ServerPort    = "8081"
	ServerAddress = "192.168.15.8"
)

func main() {
	if err := app.NewApp(ServerAddress, ServerPort).Run(); err != nil {
		log.WithField("cause", err.Error()).Fatal("error starting app")
	}
}
