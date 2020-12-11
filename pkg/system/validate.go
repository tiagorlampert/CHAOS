package system

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/pkg/util"
)

func ValidateOS() {
	osType := DetectOS()
	switch osType {
	case Linux:
		log.Info("Linux detected!")
	case Darwin:
		log.Warn("MacOS is only supported to compile to itself...")
		util.Sleep(5)
	case Windows:
		log.Fatal("Windows are not supported like Host!")
	default:
		log.Fatal("OS not supported like Host!")
	}
}
