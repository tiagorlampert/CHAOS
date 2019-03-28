package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"

	c "github.com/tiagorlampert/CHAOS/src/color"
)

const (
	NEW_LINE   string = "\n"
	BEGIN_NAME string = "chaos > "
)

func DetectOS() {
	if runtime.GOOS == "linux" {
		fmt.Println("[i] Linux!")
	} else if runtime.GOOS == "darwin" {
		ClearScreen()
		fmt.Println("[i] MacOS is only supported to compile to itself...")
		WaitTime(4)
	} else if runtime.GOOS == "windows" {
		fmt.Println("[!] Windows are not supported like Host!")
		os.Exit(0)
	} else {
		fmt.Println("[!] OS not supported like Host!")
		os.Exit(0)
	}
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func WaitTime(sec time.Duration) {
	go func() {
		time.Sleep(time.Second * sec)
	}()
	select {
	case <-time.After(time.Second * sec):
	}
}

func ReadLine() string {
	buf := bufio.NewReader(os.Stdin)
	lin, _, err := buf.ReadLine()
	if err != nil {
		fmt.Println(c.RED, "[!] Error to Read Line!")
	}
	return string(lin)
}

func RunServe() {
	go exec.Command("sh", "-c", "xterm -e \"go run src/serve/serve.go\"").Output()
}

func PauseAwaitKeyPressed() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func RemoveNewLineCharFromConnection(conn net.Conn) {
	newLineChar, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(newLineChar)
}
