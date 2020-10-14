package util

import "runtime"

type Type int

const (
	Unknown Type = iota
	Windows
	Linux
	Darwin
)

// DetectOS return an int which represent an OS type
func DetectOS() Type {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		return Darwin
	default:
		return Unknown
	}
}
