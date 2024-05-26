package main

import (
	_ "embed"
	"github.com/tiagorlampert/CHAOS/client/app"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/ui"
	"github.com/tiagorlampert/CHAOS/client/app/utils"
)

var (
	Version = "dev"
)

//go:embed config.json
var configFile []byte

func main() {
	config := utils.ReadConfigFile(configFile)

	ui.ShowMenu(Version, config.ServerAddress, config.Port)

	app.New(environment.Load(config.ServerAddress, config.Port, config.Token)).Run()
}
