package main

import (
	"github.com/tiagorlampert/CHAOS/client/app"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"github.com/tiagorlampert/CHAOS/client/app/ui"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
)

var (
	Version       = "dev "
	ServerPort    = ""
	ServerAddress = ""
	Token         = ""
)

func main() {
	ui.ShowLogo(Version)
	appConfiguration := environment.LoadConfiguration(ServerAddress, ServerPort, Token)

	httpClient := network.NewHttpClient(10)
	app.NewApp(httpClient, appConfiguration).Run()
}
