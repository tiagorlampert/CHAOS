package system

import (
	"os"
	"os/exec"
	"runtime"
)

type OSType int

const (
	Unknown OSType = iota
	Windows
	Linux
	//Darwin
)

var OSTargetMap = map[OSType]string{
	Windows: "Windows",
	Linux:   "Linux",
	//Darwin:  "Mac OS",
}

var OSTargetIntMap = map[int]OSType{
	1: Windows,
	2: Linux,
	//3:  Darwin,
}

// DetectOS return an int which represent a OS type
func DetectOS() OSType {
	switch runtime.GOOS {
	case `windows`:
		return Windows
	case `linux`:
		return Linux
	//case `darwin`:
	//	return Darwin
	default:
		return Unknown
	}
}

// ClearScreen clean the terminal
func ClearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
