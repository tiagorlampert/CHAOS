package os

import (
	"runtime"
)

type OSType int

const (
	Unknown OSType = iota
	Windows
	Linux
	//Darwin
)

var TargetMap = map[OSType]string{
	Windows: "Windows",
	Linux:   "Linux",
	//Darwin:  "Mac OS",
}

var TargetIntMap = map[int]OSType{
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
