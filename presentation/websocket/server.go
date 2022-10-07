package websocket

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"net/http"
)

func NewServer(configuration *environment.Configuration) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", configuration.Server.Port), nil)
}
