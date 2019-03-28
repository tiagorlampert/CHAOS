package network

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"

	c "github.com/tiagorlampert/CHAOS/src/color"
	"github.com/tiagorlampert/CHAOS/src/util"
)

func SendMessage(conn net.Conn, message string) {
	conn.Write([]byte(base64.URLEncoding.EncodeToString([]byte(message)) + util.NEW_LINE))
}

func SendMessageByte(conn net.Conn, message []byte) {
	conn.Write([]byte(base64.URLEncoding.EncodeToString([]byte(message)) + util.NEW_LINE))
}

func SendMessageRaw(conn net.Conn, message string) {
	conn.Write([]byte(message + util.NEW_LINE))
}

func ReceiveMessage(conn net.Conn) {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageDecoded, _ := base64.StdEncoding.DecodeString(message)
	fmt.Println(c.YELLOW, string(messageDecoded)+"\n")
}

func ReceiveMessageReturn(conn net.Conn) string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	return string(message)
}

func ReceiveMessageReturnDecodeString(conn net.Conn) string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageDecoded, _ := base64.StdEncoding.DecodeString(message)
	return string(messageDecoded)
}
