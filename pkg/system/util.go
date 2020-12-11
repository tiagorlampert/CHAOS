package system

import (
	"fmt"
	"os"
	"os/exec"
)

func ClearScreen() {
	osType := DetectOS()
	switch osType {
	case Linux:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func QuitApp() {
	ClearScreen()
	fmt.Println("Bye, See you later!")
	os.Exit(0)
}
