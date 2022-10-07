package delete

import (
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"os"
)

type Service struct {
}

func NewService() services.Delete {
	return &Service{}
}

func (d Service) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
