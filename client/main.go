package main

import (
	"github.com/tiagorlampert/CHAOS/client/app"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/ui"
)

var (
	Version       = "dev"
	Port          = ""
	ServerAddress = ""
	Token         = ""
)

func main() {
	ui.ShowMenu(Version, ServerAddress, Port)

	app.New(environment.Load(ServerAddress, Port, Token)).Run()
}
