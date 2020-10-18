package download

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"github.com/tiagorlampert/CHAOS/client/app/util/os"
	"net"
)

type DownloadUseCase struct {
	Connection net.Conn
}

func NewDownloadUseCase(conn net.Conn) usecase.Download {
	return &DownloadUseCase{
		Connection: conn,
	}
}

func (d DownloadUseCase) File() {
	_ = network.Send(d.Connection, []byte("ok"))

	// read file path
	path, _ := network.Read(d.Connection)

	file, err := os.ReadFile(string(path))
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading file")
	}

	if err = network.Send(d.Connection, file); err != nil {
		log.WithField("cause", err.Error()).Error("error sending file")
	}
}
