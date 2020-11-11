package util

import (
	"bufio"
	"os"
)

func EnterAnyKey() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
